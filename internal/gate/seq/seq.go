package seq

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"sync"
	"time"
)

type Sequence interface {
	Next(ctx context.Context, key string) (uint64, error)
}

type RedisSequence struct {
	redisClient *redis.Client

	mu         sync.Mutex
	bufferSize int64
	buffers    map[string][]uint64
	nextIndex  map[string]int
}

const (
	cacheSeqPrefix = "seq:%s"
	cacheSeqLock   = "seq:lock:%s"
	lockExpireTime = 3 * time.Second
)

func NewRedisSequence(redisClient *redis.Client, bufferSize int64) Sequence {
	return &RedisSequence{
		redisClient: redisClient,
		bufferSize:  bufferSize,
		buffers:     make(map[string][]uint64),
		nextIndex:   make(map[string]int),
	}
}

func (s *RedisSequence) Next(ctx context.Context, key string) (uint64, error) {
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

func (s *RedisSequence) fillBufferWithLock(ctx context.Context, key string) error {
	var unlock func()
	for attempt := 0; attempt < 3; attempt++ {
		lock, err := s.redisClient.SetNX(ctx, fmt.Sprintf(cacheSeqLock, key), 1, lockExpireTime).Result()
		if err != nil {
			return err
		}
		if lock {
			unlock = func() { s.redisClient.Del(ctx, fmt.Sprintf(cacheSeqLock, key)) }
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

func (s *RedisSequence) fillBuffer(ctx context.Context, key string) error {
	cacheKey := fmt.Sprintf(cacheSeqPrefix, key)
	current, err := s.redisClient.Get(ctx, cacheKey).Int64()
	if err != nil && !errors.Is(err, redis.Nil) {
		return err
	}

	// todo redis查不到时，尝试从db中获取key的当前最大值

	// 设置新的序列号范围
	start := current + 1
	end := start + s.bufferSize - 1

	// 使用 Redis 的 IncrBy 命令获取新的序列号范围
	newCurrent, err := s.redisClient.IncrBy(ctx, cacheKey, s.bufferSize).Result()
	if err != nil {
		return err
	}
	if newCurrent != end {
		return errors.New("sequence number mismatch")
	}

	// 重置缓冲区
	s.buffers[key] = make([]uint64, s.bufferSize)
	for i := 0; i < int(s.bufferSize); i++ {
		s.buffers[key][i] = uint64(start + int64(i))
	}
	s.nextIndex[key] = 0
	return nil
}
