package seq

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"ppim/internal/logic/store"
	"sync"
	"time"
)

// Sequence 序列号生成器可抽离为单独的服务，避免每个logic服务节点维护管理自己的序列号
type Sequence interface {
	Next(ctx context.Context, key string) (uint64, error)
}

type RedisSequence struct {
	redisClient *redis.Client

	mu         sync.Mutex
	bufferSize uint64
	buffers    map[string][]uint64
	nextIndex  map[string]int
}

const (
	cacheSeqPrefix = "seq:%s"
	cacheSeqLock   = "seq:lock:%s"
	lockExpireTime = 3 * time.Second
)

func NewRedisSequence(redisClient *redis.Client, bufferSize uint64) Sequence {
	return &RedisSequence{
		redisClient: redisClient,
		bufferSize:  bufferSize,
		buffers:     make(map[string][]uint64),
		nextIndex:   make(map[string]int),
	}
}

func (s *RedisSequence) Next(ctx context.Context, key string) (uint64, error) {
	// 使用incr时，利用redis实现全局序列号，但高并发时redis压力大；
	// 使用incrWithBuf时，需要保证同一会话key路由至同一logic节点，且扩缩容时重置每个节点的缓冲区
	return s.incrWithBuf(ctx, key)
}

// redis incr生成序列号
func (s *RedisSequence) incr(ctx context.Context, key string) (uint64, error) {
	cacheKey := fmt.Sprintf(cacheSeqPrefix, key)
	current, err := s.redisClient.Incr(ctx, cacheKey).Uint64()
	if err != nil {
		return 0, err
	}

	if current == 1 { // 获取序列号=1时可能缓存丢失，尝试从db中获取key当前最大值初始化缓存
		current, err = new(store.Conversation).GetMaxSeq(ctx, key)
		if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
			return 0, err
		}
		current, err = s.redisClient.IncrBy(ctx, cacheKey, int64(current)+1).Uint64()
	}
	return current, nil
}

// redis incrBy + 缓冲区生成序列号
func (s *RedisSequence) incrWithBuf(ctx context.Context, key string) (uint64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 缓冲区为空，从 Redis 获取新的序列号
	if s.nextIndex[key] >= len(s.buffers[key]) {
		err := s.fillBufferWithLock(ctx, key)
		if err != nil {
			return 0, err
		}
	}

	// 从缓冲区中获取下一个序列号
	val := s.buffers[key][s.nextIndex[key]]
	s.nextIndex[key]++
	return val, nil
}

func (s *RedisSequence) fillBuffer(ctx context.Context, key string) error {
	cacheKey := fmt.Sprintf(cacheSeqPrefix, key)
	current, err := s.redisClient.Get(ctx, cacheKey).Uint64()
	if err != nil && !errors.Is(err, redis.Nil) {
		return err
	}

	if current == 0 { // 未获取到序列号时可能缓存丢失，尝试从db中获取key当前最大值初始化缓存
		current, err = new(store.Conversation).GetMaxSeq(ctx, key)
		if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
			return err
		}
	}

	// 设置新的序列号范围
	start := current + 1
	end := start + s.bufferSize - 1

	newCurrent, err := s.redisClient.IncrBy(ctx, cacheKey, int64(s.bufferSize)).Uint64()
	if err != nil {
		return err
	}
	if newCurrent != end {
		return errors.New("sequence number mismatch")
	}

	// 重置缓冲区
	s.buffers[key] = make([]uint64, s.bufferSize)
	for i := range s.bufferSize {
		s.buffers[key][i] = start + i
	}
	s.nextIndex[key] = 0
	return nil
}

func (s *RedisSequence) fillBufferWithLock(ctx context.Context, key string) error {
	var (
		lockKey = fmt.Sprintf(cacheSeqLock, key)
		unlock  func()
	)
	for attempt := 0; attempt < 3; attempt++ {
		lock, err := s.redisClient.SetNX(ctx, lockKey, 1, lockExpireTime).Result()
		if err != nil {
			return err
		}
		if lock {
			unlock = func() { s.redisClient.Del(ctx, lockKey) }
			break
		}
		// 等待一段时间后重试
		select {
		case <-time.After(300 * time.Millisecond):
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	if unlock == nil {
		return errors.New("failed to acquire lock")
	}
	defer unlock()

	return s.fillBuffer(ctx, key)
}
