package net

import (
	"context"
	"sync"
	"time"
)

type RetryMsg struct {
	MsgBytes []byte // 消息数据
	ConnFD   int    // 连接ID
	MsgID    string // 消息ID
	UID      string // 用户ID
	retry    int    // 重试次数
}

type RetryQueue struct {
	mtx      sync.RWMutex
	elements []*RetryMsg
	used     bool
}

func (q *RetryQueue) Enqueue(ele *RetryMsg) {
	q.mtx.Lock()
	defer q.mtx.Unlock()

	q.elements = append(q.elements, ele)
}

func (q *RetryQueue) Dequeue() interface{} {
	q.mtx.Lock()
	defer q.mtx.Unlock()

	if len(q.elements) == 0 {
		return nil
	}
	element := q.elements[0]
	q.elements = q.elements[1:]
	return element
}

func (q *RetryQueue) Take(num int) []*RetryMsg {
	q.mtx.Lock()
	defer q.mtx.Unlock()

	if len(q.elements) < num {
		num = len(q.elements)
	}
	if num == 0 {
		return nil
	}
	els := q.elements[:num]
	q.elements = q.elements[num:]
	return els
}

func (q *RetryQueue) Size() int {
	q.mtx.RLock()
	defer q.mtx.RUnlock()

	return len(q.elements)
}

// RetryDelivery 通过重试队列保障消息可靠投递
type RetryDelivery struct {
	svc      *ServerContext
	queue    *RetryQueue
	maxRetry int

	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
}

func newRetryDelivery(svc *ServerContext) *RetryDelivery {
	ctx, cancel := context.WithCancel(context.Background())
	return &RetryDelivery{
		svc:      svc,
		ctx:      ctx,
		cancel:   cancel,
		maxRetry: 3,
		queue: &RetryQueue{
			elements: make([]*RetryMsg, 0, 1024),
		},
	}
}

func (r *RetryDelivery) Add(msg *RetryMsg) {
	r.queue.Enqueue(msg)
}

func (r *RetryDelivery) Remove(msgId string) {
	// todo
}

func (r *RetryDelivery) Start() {
	r.wg.Add(1)
	go r.work()
}

func (r *RetryDelivery) Stop() {
	r.cancel()
	r.queue = nil
	r.wg.Wait()
}

func (r *RetryDelivery) work() {
	defer r.wg.Done()

	ticker := time.NewTicker(3 * time.Second)
	for {
		select {
		case <-r.ctx.Done():
			return
		case <-ticker.C:
			if r.queue.used {
				return
			}

			func() {
				r.queue.used = true
				defer func() {
					r.queue.used = false
				}()

				ms := r.queue.Take(1000)
				if len(ms) == 0 {
					return
				}
				for _, m := range ms {
					m.retry++

					client := r.svc.ConnManager.GetWithFD(m.ConnFD)
					if client == nil {
						return
					}

					_, _ = client.Write(m.MsgBytes)

					if m.retry < r.maxRetry {
						r.Add(m)
					}
				}
			}()
		}
	}
}
