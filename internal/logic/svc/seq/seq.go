package seq

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"ppim/internal/logic/store"
	"ppim/pkg/ext"
	"ppim/pkg/util"
	"sync"
)

// Sequence 序列号生成器，生成key对应的下一序列号
type Sequence interface {
	Next(ctx context.Context, key string) (uint64, error)
}

// SequenceStorage 序列号存储器，存储key对应的最大序列号
type SequenceStorage interface {
	GetMaxSeq(ctx context.Context, key string) (uint64, error)
	SetMaxSeq(ctx context.Context, key string, maxSeq uint64) error
}

const (
	cacheSeqPrefix = "seq:%s"
)

type StepSequence struct {
	step    uint64
	curSeq  sync.Map         // key的当前序列号
	maxSeq  sync.Map         // key的最大序列号（优化：分段共享maxSeq）
	locker  *ext.SegmentLock // 分段锁
	storage SequenceStorage
}

func NewStepSequence(step uint64, store SequenceStorage) Sequence {
	return &StepSequence{
		step:    step,
		storage: store,
		locker:  ext.NewSegmentLock(10),
	}
}

func (s *StepSequence) Next(ctx context.Context, key string) (uint64, error) {
	index := int(util.Murmur32(key))
	s.locker.Lock(index)
	defer s.locker.Unlock(index)

	var curSeq, maxSeq uint64 = 0, 0
	if curVal, ok := s.curSeq.Load(key); ok {
		curSeq = curVal.(uint64)
	}
	if maxVal, ok := s.maxSeq.Load(key); ok { // 优化：分段共享时，此处可按key的进行hash获取
		maxSeq = maxVal.(uint64)
	}

	curSeq++
	if curSeq > maxSeq {
		if maxSeq == 0 {
			curMaxSeq, err := s.storage.GetMaxSeq(ctx, key)
			if err != nil {
				return 0, err
			}
			maxSeq = curMaxSeq
			curSeq = curMaxSeq + 1
		}
		newMaxSeq := maxSeq + s.step
		err := s.storage.SetMaxSeq(ctx, key, newMaxSeq)
		if err != nil {
			return 0, fmt.Errorf("failed to save max_seq: %v", err)
		}
		s.maxSeq.Store(key, newMaxSeq)
	}
	s.curSeq.Store(key, curSeq)
	return curSeq, nil
}

type DefaultSeqStorage struct {
	cache    *redis.Client
	seqStore *store.Conversation
}

func WithDefaultSeqStorage(cache *redis.Client, dbStore *store.Conversation) SequenceStorage {
	return &DefaultSeqStorage{
		cache:    cache,
		seqStore: dbStore,
	}
}

func (s *DefaultSeqStorage) GetMaxSeq(ctx context.Context, key string) (uint64, error) {
	curMaxSeq, err := s.cache.Get(ctx, fmt.Sprintf(cacheSeqPrefix, key)).Uint64()
	if err != nil && !errors.Is(err, redis.Nil) {
		return 0, fmt.Errorf("failed to get max_seq: %v", err)
	}
	// 未获取到最大序列号可能缓存丢失，尝试从db中获取
	if curMaxSeq == 0 {
		curMaxSeq, err = new(store.Conversation).GetMaxSeq(ctx, key)
		if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
			return 0, fmt.Errorf("failed to get max_seq: %v", err)
		}
	}
	return curMaxSeq, nil
}

func (s *DefaultSeqStorage) SetMaxSeq(ctx context.Context, key string, maxSeq uint64) error {
	err := s.cache.Set(ctx, fmt.Sprintf(cacheSeqPrefix, key), maxSeq, 0).Err()
	if err != nil {
		return fmt.Errorf("failed to save max_seq: %v", err)
	}
	return nil
}

type RedisSequence struct {
	store *redis.Client
}

func NewRedisSequence(client *redis.Client) Sequence {
	return &RedisSequence{
		store: client,
	}
}

func (s *RedisSequence) Next(ctx context.Context, key string) (uint64, error) {
	cacheKey := fmt.Sprintf(cacheSeqPrefix, key)
	current, err := s.store.Incr(ctx, cacheKey).Uint64()
	if err != nil {
		return 0, err
	}

	if current == 1 { // 获取序列号=1时可能缓存丢失，尝试从db中获取key当前最大值初始化缓存
		current, err = new(store.Conversation).GetMaxSeq(ctx, key)
		if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
			return 0, err
		}
		current, err = s.store.IncrBy(ctx, cacheKey, int64(current)+1).Uint64()
	}
	return current, nil
}

//func (s *RedisSequence) withRetryLock(ctx context.Context, key string, fn func() error) error {
//	var (
//		lockKey = fmt.Sprintf(cacheSeqLock, key)
//		unlock  func()
//	)
//	for attempt := 0; attempt < 3; attempt++ {
//		locker, err := s.storage.SetNX(ctx, lockKey, 1, lockExpireTime).Result()
//		if err != nil {
//			return err
//		}
//		if locker {
//			unlock = func() { s.storage.Del(ctx, lockKey) }
//			break
//		}
//		// 等待一段时间后重试
//		select {
//		case <-time.After(200 * time.Millisecond)
//		case <-ctx.Done():
//			return ctx.Err()
//		}
//	}
//	if unlock == nil {
//		return errors.New("failed to acquire locker")
//	}
//	defer unlock()
//
//	return fn()
//}
