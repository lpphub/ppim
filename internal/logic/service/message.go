package service

import (
	"context"
	"ppim/internal/chat"
	"ppim/internal/logic/global"
	"ppim/internal/logic/store"
	"ppim/internal/logic/types"
	"time"
)

type MessageSrv struct{}

var MsgSrv = &MessageSrv{}

func (s *MessageSrv) HandleMsg(ctx context.Context, msg *types.MessageDTO) error {
	// 1. 消息持久化
	mm := &store.Message{
		MsgID:            msg.MsgID,
		MsgSeq:           msg.MsgSeq,
		MsgNo:            msg.MsgNo,
		MsgType:          msg.MsgType,
		Content:          msg.Content,
		ConversationID:   msg.ConversationID,
		ConversationType: msg.ConversationType,
		FromID:           msg.FromID,
		ToID:             msg.ToID,
		CreatedAt:        time.Now(),
	}
	if err := mm.Insert(ctx); err != nil {
		return err
	}

	// 2. 消息推送
	var receivers []string
	if msg.ConversationType == chat.ConvSingle {
		receivers = append(receivers, msg.ToID)
	} else if msg.ConversationType == chat.ConvGroup {
		members, err := new(store.Group).ListMembers(ctx, msg.ToID)
		if err != nil {
			return err
		}
		receivers = append(receivers, members...)
	}
	if len(receivers) == 0 {
		return nil
	}
	// todo 写扩散 索引每个接收者的最近会话
	if err := IndexConversation(ctx, receivers, msg); err != nil {
		return err
	}

	var (
		onlineUserDeviceSlice []string
		offlineUIDSlice       []string
	)
	for _, uid := range receivers {
		online, _ := global.Redis.SMembers(ctx, OnSrv.getOnlineKey(uid)).Result()
		if len(online) > 0 {
			onlineUserDeviceSlice = append(onlineUserDeviceSlice, online...)
		} else {
			offlineUIDSlice = append(offlineUIDSlice, uid)
		}
	}
	if len(onlineUserDeviceSlice) > 0 {
		// todo 推送消息（在线）
	}
	if len(offlineUIDSlice) > 0 {
		// todo 消息通知（离线）
	}
	return nil
}
