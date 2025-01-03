package net

import (
	"context"
	"fmt"
	"github.com/lpphub/golib/logger"
	"ppim/pkg/util"
	"sync"
	"sync/atomic"
	"time"
)

type RetryMsg struct {
	MsgBytes []byte // 消息数据
	ConnFD   int    // 连接ID
	MsgID    string // 消息ID
	UID      string // 用户ID
	retry    int    // 重试次数
	deleted  bool   // 是否删除
}

type RetryDelivery struct {
	queue []*RetryMsg
	hash  map[string]*RetryMsg

	svc        *ServerContext
	maxRetries int

	working atomic.Bool
	mtx     sync.Mutex
	ctx     context.Context
	cancel  context.CancelFunc
}

func newRetryDelivery(svc *ServerContext, maxRetries int) *RetryDelivery {
	ctx, cancel := context.WithCancel(context.Background())
	return &RetryDelivery{
		svc:        svc,
		maxRetries: maxRetries,
		ctx:        ctx,
		cancel:     cancel,
		queue:      make([]*RetryMsg, 0, 2048),
		hash:       make(map[string]*RetryMsg),
	}
}

func (r *RetryDelivery) Add(el *RetryMsg) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	r.queue = append(r.queue, el)
	r.hash[r.getUk(el.ConnFD, el.MsgID)] = el
}

func (r *RetryDelivery) Remove(connFD int, msgId string) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	key := r.getUk(connFD, msgId)
	if el, ok := r.hash[key]; ok {
		el.deleted = true
		delete(r.hash, key)
	}
}

func (r *RetryDelivery) getUk(connFD int, msgId string) string {
	return fmt.Sprintf("%d:%s", connFD, msgId)
}

func (r *RetryDelivery) Take(num int) []*RetryMsg {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	if len(r.queue) == 0 {
		return nil
	}
	if num > len(r.queue) {
		num = len(r.queue)
	}

	batch := r.queue[:num]
	r.queue = r.queue[num:]
	return batch
}

func (r *RetryDelivery) start() {
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		for {
			select {
			case <-r.ctx.Done():
				return
			case <-ticker.C:
				logger.Log().Debug().Msgf("重试消息 size: %d", len(r.queue))

				go util.WithRecover(r.work)
			}
		}
	}()
}

func (r *RetryDelivery) stop() {
	r.cancel()
}

func (r *RetryDelivery) work() {
	if r.working.Load() {
		return
	}
	r.working.Store(true)
	defer r.working.Store(false)

	for {
		select {
		case <-r.ctx.Done():
			return
		default:
			ms := r.Take(100)
			if len(ms) == 0 {
				return
			}
			for _, m := range ms {
				r.runRetry(m)
			}
		}
	}
}

func (r *RetryDelivery) runRetry(msg *RetryMsg) {
	if msg.deleted {
		return
	}
	msg.retry++
	logger.Log().Info().Msgf("重试消息发送：uid=%s, msgID=%s, retryCount=%d", msg.UID, msg.MsgID, msg.retry)

	client := r.svc.ConnManager.GetWithFD(msg.ConnFD)
	if client == nil {
		return
	}

	_, _ = client.Write(msg.MsgBytes)

	if msg.retry >= r.maxRetries {
		r.Remove(msg.ConnFD, msg.MsgID)
	} else {
		r.Add(msg)
	}
}

// RetryManager 通过重试队列保障消息可靠投递
type RetryManager struct {
	svc     *ServerContext
	size    int
	buckets []*RetryDelivery
}

func newRetryManager(svc *ServerContext, size int) *RetryManager {
	return &RetryManager{
		svc:     svc,
		size:    size,
		buckets: make([]*RetryDelivery, size),
	}
}

func (r *RetryManager) Add(msg *RetryMsg) {
	index := (int(util.Murmur32(msg.MsgID))%r.size + r.size) % r.size
	r.buckets[index].Add(msg)
}

func (r *RetryManager) Remove(connFD int, msgId string) {
	index := (int(util.Murmur32(msgId))%r.size + r.size) % r.size
	r.buckets[index].Remove(connFD, msgId)
}

func (r *RetryManager) Start() {
	for i := 0; i < r.size; i++ {
		r.buckets[i] = newRetryDelivery(r.svc, 1)
		r.buckets[i].start()
	}
}

func (r *RetryManager) Stop() {
	for _, t := range r.buckets {
		t.stop()
	}
}
