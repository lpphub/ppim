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

// Sequence 序列号生成器可抽离为单独的服务，避免每个logic服务节点维护管理自己的序列号
type Sequence interface {
	Next(ctx context.Context, key string) (uint64, error)
}

const (
	cacheSeqPrefix = "seq:%s"
)

type StepSequence struct {
	step   uint64
	curSeq sync.Map         // key的当前序列号
	maxSeq sync.Map         // key的最大序列号（可优化多个key共享一个maxSeq以减少IO操作）
	locker *ext.SegmentLock // 分段锁
	store  *redis.Client
}

func NewStepSequence(store *redis.Client, step uint64) Sequence {
	return &StepSequence{
		step:   step,
		store:  store,
		locker: ext.NewSegmentLock(10),
	}
}

func (s *StepSequence) Next(ctx context.Context, key string) (uint64, error) {
	index := int(util.Murmur32(key))
	s.locker.Lock(index)
	defer s.locker.Unlock(index)

	var curSeq, maxSeq uint64 = 0, 0
	if curr, ok := s.curSeq.Load(key); ok {
		curSeq = curr.(uint64)
	}
	if maxs, ok := s.maxSeq.Load(key); ok {
		maxSeq = maxs.(uint64)
	}

	curSeq++
	if curSeq > maxSeq {
		storeKey := fmt.Sprintf(cacheSeqPrefix, key)
		if maxSeq == 0 {
			curMaxSeq, err := s.store.Get(ctx, storeKey).Uint64()
			if err != nil && !errors.Is(err, redis.Nil) {
				return 0, fmt.Errorf("failed to get max_seq: %v", err)
			}
			if curMaxSeq == 0 { // 未获取到最大序列号可能缓存丢失，尝试从db中获取
				curMaxSeq, err = new(store.Conversation).GetMaxSeq(ctx, key)
				if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
					return 0, fmt.Errorf("failed to get max_seq: %v", err)
				}
			}
			maxSeq = curMaxSeq
			curSeq = curMaxSeq + 1
		}
		newMaxSeq := maxSeq + s.step
		err := s.store.Set(ctx, storeKey, newMaxSeq, 0).Err()
		if err != nil {
			return 0, fmt.Errorf("failed to save max_seq: %v", err)
		}
		s.maxSeq.Store(key, newMaxSeq)
	}

	s.curSeq.Store(key, curSeq)
	return curSeq, nil
}

type GlobalSequence struct {
	store *redis.Client
}

func NewGlobalSequence(redisClient *redis.Client) Sequence {
	return &GlobalSequence{
		store: redisClient,
	}
}

func (s *GlobalSequence) Next(ctx context.Context, key string) (uint64, error) {
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

//func (s *GlobalSequence) withRetryLock(ctx context.Context, key string, fn func() error) error {
//	var (
//		lockKey = fmt.Sprintf(cacheSeqLock, key)
//		unlock  func()
//	)
//	for attempt := 0; attempt < 3; attempt++ {
//		locker, err := s.store.SetNX(ctx, lockKey, 1, lockExpireTime).Result()
//		if err != nil {
//			return err
//		}
//		if locker {
//			unlock = func() { s.store.Del(ctx, lockKey) }
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
