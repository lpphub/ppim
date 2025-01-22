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
	Retries  int    // 重试次数
}

type RetryDelivery struct {
	queue []*RetryMsg
	index map[string]int

	svc        *ServerContext
	maxRetries int

	running atomic.Bool
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
		index:      make(map[string]int),
	}
}

func (r *RetryDelivery) Add(el *RetryMsg) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	if len(r.queue) >= cap(r.queue) {
		logger.Log().Warn().Msg("重试队列已满，消息将被丢弃")
		return
	}

	r.queue = append(r.queue, el)
	r.index[r.genUK(el.ConnFD, el.MsgID)] = len(r.queue) - 1
}

func (r *RetryDelivery) Remove(connFD int, msgId string) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	key := r.genUK(connFD, msgId)
	idx, ok := r.index[key]
	if !ok {
		return
	}

	// 将要删除的元素与最后一个元素交换
	lastIdx := len(r.queue) - 1
	lastEle := r.queue[lastIdx]
	r.queue[idx] = lastEle
	r.index[r.genUK(lastEle.ConnFD, lastEle.MsgID)] = idx

	// 删除最后一个元素
	r.queue = r.queue[:lastIdx]
	delete(r.index, key)
}

func (r *RetryDelivery) genUK(connFD int, msgId string) string {
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
		ticker := time.NewTicker(3 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-r.ctx.Done():
				return
			case <-ticker.C:
				if !r.running.Load() {
					logger.Log().Debug().Msgf("重试消息 size: %d", len(r.queue))

					r.run()
				}
			}
		}
	}()
}

func (r *RetryDelivery) stop() {
	r.cancel()
}

func (r *RetryDelivery) run() {
	r.running.Store(true)
	defer r.running.Store(false)

	ms := r.Take(1000)
	if len(ms) == 0 {
		return
	}
	for _, m := range ms {
		r.runRetry(m)
	}
}

func (r *RetryDelivery) runRetry(msg *RetryMsg) {
	msg.Retries++
	logger.Log().Info().Msgf("重试消息发送：uid=%s, msgID=%s, retryCount=%d", msg.UID, msg.MsgID, msg.Retries)

	client := r.svc.ConnManager.GetWithFD(msg.ConnFD)
	if client == nil {
		return
	}

	_, _ = client.Write(msg.MsgBytes)

	if msg.Retries >= r.maxRetries {
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
	rm := &RetryManager{
		svc:     svc,
		size:    size,
		buckets: make([]*RetryDelivery, size),
	}

	// 启动
	rm.Start()
	return rm
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
