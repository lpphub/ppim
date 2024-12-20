package task

import (
	"context"
	"ppim/internal/gate/net"
	"sync"
)

type RetryMsg struct {
	MsgBytes []byte // 消息数据
	ConnFD   int    // 连接ID
	MsgID    string // 消息ID
	UID      string // 用户ID
	retry    int    // 重试次数
}

// RetryDelivery 通过重试队列保障消息可靠投递
type RetryDelivery struct {
	svc      *net.ServerContext
	queue    chan *RetryMsg
	workers  int
	maxRetry int

	//mtx    sync.RWMutex
	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
}

func NewRetryDelivery(svc *net.ServerContext) *RetryDelivery {
	ctx, cancel := context.WithCancel(context.Background())
	return &RetryDelivery{
		svc:    svc,
		queue:  make(chan *RetryMsg, 1024),
		ctx:    ctx,
		cancel: cancel,
	}
}

func (r *RetryDelivery) Add(msg *RetryMsg) {
	select {
	case <-r.ctx.Done():
		return
	case r.queue <- msg:
	}
}

func (r *RetryDelivery) Remove(msgId string) {
	// todo
}

func (r *RetryDelivery) Start() {
	if r.workers <= 0 {
		r.workers = 1
	}
	for range r.workers {
		r.wg.Add(1)
		go r.work()
	}
}

func (r *RetryDelivery) Stop() {
	r.cancel()
	close(r.queue)
	r.wg.Wait()
}

func (r *RetryDelivery) work() {
	defer r.wg.Done()

	for {
		select {
		case <-r.ctx.Done():
			return
		case t, ok := <-r.queue:
			if !ok {
				return
			}
			t.retry++

			client := r.svc.ConnManager.GetWithFD(t.ConnFD)
			if client == nil {
				return
			}

			_, _ = client.Write(t.MsgBytes)

			if t.retry < r.maxRetry {
				r.Add(t)
			}
		}
	}
}
