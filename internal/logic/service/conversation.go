package service

import (
	"context"
	"hash/adler32"
	"ppim/internal/logic/store"
	"ppim/internal/logic/types"
	"ppim/pkg/ext"
	"time"
)

// ConversationSrv
/**
 * 写扩散：每个用户对应一个timeline, 消息到达后每个接收者更新自身timeline
 * 读扩散：一个会话对应一个timeline，消息到达后更新此会话最新timeline
 */
type ConversationSrv struct {
	segmentLock ext.SegmentLock
}

const (
	cacheConvRecent = "conv:recent:%s"
)

func (c *ConversationSrv) IndexConv(ctx context.Context, msg *types.MessageDTO, uidSlice []string) error {
	// 发送者会话
	uidSlice = append(uidSlice, msg.FromID)

	// 接收者会话
	for _, uid := range uidSlice {
		c.indexWithLock(ctx, msg, uid)
	}
	return nil
}

func (c *ConversationSrv) indexWithLock(ctx context.Context, msg *types.MessageDTO, uid string) {
	index := adler32.Checksum([]byte(uid))
	c.segmentLock.Lock(index)
	defer c.segmentLock.Unlock(index)

	conv, _ := new(store.Conversation).GetOne(ctx, uid, msg.ConversationID)
	if conv == nil {
		conv = &store.Conversation{
			ConversationID:   msg.ConversationID,
			ConversationType: msg.ConversationType,
			UID:              uid,
			UnreadCount:      1,
			LastMsgId:        msg.MsgID,
			LastMsgSeq:       msg.MsgSeq,
			CreatedAt:        time.Now(),
		}
		_ = conv.Insert(ctx)
	} else {
		conv.UnreadCount++
		conv.LastMsgId = msg.MsgID
		conv.LastMsgSeq = msg.MsgSeq
		conv.UpdatedAt = time.Now()
		_ = conv.Update(ctx)
	}

}
