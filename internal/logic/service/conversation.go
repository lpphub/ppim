package service

import (
	"context"
	"ppim/internal/logic/store"
	"ppim/internal/logic/types"
	"time"
)

// A与B聊天，A与C聊天，A与群聊D聊天
// conv:recent:{uid} single|A@B 1, single|A@C 2, group|{groupId} 3

type ConversationSrv struct{}

const (
	cacheConvRecent   = "conv:recent:%s"
	cacheConvTimeline = "conv:timeline:%s"
)

func IndexConversation(ctx context.Context, uidSlice []string, msg *types.MessageDTO) error {
	// 发送者会话
	uidSlice = append(uidSlice, msg.FromID)

	// 接收者会话
	for _, uid := range uidSlice {
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
	return nil
}
