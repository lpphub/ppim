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
	var (
		onlineReceivers []string
		offlineUidList  []string
	)
	if msg.ConversationType == chat.ConvSingle {
		online, _ := global.Redis.SMembers(ctx, OnSrv.getOnlineKey(msg.ToID)).Result()
		if len(online) > 0 { // 用户所有在线设备
			onlineReceivers = append(onlineReceivers, online...)
		} else {
			offlineUidList = append(offlineUidList, msg.ToID)
		}
	} else if msg.ConversationType == chat.ConvGroup {
		members, err := new(store.Group).ListMembers(ctx, msg.ToID)
		if err != nil {
			return err
		}
		// 获取在线成员
		for _, member := range members {
			online, _ := global.Redis.SMembers(ctx, OnSrv.getOnlineKey(member)).Result()
			if len(online) > 0 {
				onlineReceivers = append(onlineReceivers, online...)
			} else {
				offlineUidList = append(offlineUidList, member)
			}
		}
	}
	if len(onlineReceivers) > 0 {
		// todo 2.1 推送消息（在线）
	}

	if len(offlineUidList) > 0 {
		// todo 2.2 消息通知（离线）
	}
	return nil
}
