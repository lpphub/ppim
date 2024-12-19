package task

import (
	"context"
	"ppim/internal/chatlib"
	"ppim/pkg/ext"
)

// RetryDelivery 通过重试队列保障消息可靠投递
type RetryDelivery struct {
	queue *ext.Queue
	mtx   ext.SegmentLock
}

func (r *RetryDelivery) Add(msg *chatlib.DeliveryMsg) error {
	uid := chatlib.DigitizeUID(msg.FromUID)
	r.mtx.Lock(uid)
	defer r.mtx.Unlock(uid)

	r.queue.Add(r.retry)

	return nil
}

func (r *RetryDelivery) retry(ctx context.Context) error {
	return nil
}
