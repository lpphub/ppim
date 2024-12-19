package task

import (
	"context"
	"ppim/internal/chatlib"
	"ppim/internal/logic/types"
	"ppim/pkg/ext"
)

// RetryDelivery 通过重试队列保障消息可靠投递
type RetryDelivery struct {
	queue *ext.Queue
	mtx   ext.SegmentLock
}

func (r *RetryDelivery) Add(msg *types.MessageDTO) error {
	uid := chatlib.DigitizeUID(msg.FromID)
	r.mtx.Lock(uid)
	defer r.mtx.Unlock(uid)

	r.queue.Add(r.retry)

	return nil
}

func (r *RetryDelivery) retry(ctx context.Context) error {
	return nil
}
