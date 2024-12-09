package service

import (
	"context"
	"github.com/lpphub/golib/logger"
	"github.com/pkg/errors"
	"ppim/internal/chatlib"
	"ppim/internal/logic/global"
	"ppim/internal/logic/store"
	"ppim/internal/logic/types"
	"time"
)

var (
	ErrMsgStore  = errors.New("消息持久化失败")
	ErrConvIndex = errors.New("消息索引更新失败")
	ErrMsgRoute  = errors.New("消息路由失败")
)

type MessageSrv struct{}

func (s *MessageSrv) HandleMsg(ctx context.Context, msg *types.MessageDTO) error {
	ctx = logger.WithCtx(ctx)

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
		logger.Err(ctx, err, "")
		return ErrMsgStore
	}

	var receivers []string
	if msg.ConversationType == chatlib.ConvSingle {
		receivers = append(receivers, msg.ToID)
	} else if msg.ConversationType == chatlib.ConvGroup {
		members, err := new(store.Group).ListMembers(ctx, msg.ToID)
		if err != nil {
			logger.Err(ctx, err, "")
		} else {
			receivers = append(receivers, members...)
		}
	}
	if len(receivers) == 0 {
		logger.Warn(ctx, "消息接收者列表为空")
		return nil
	}
	// 索引会话
	if err := svc.ConvSrv.IndexConv(ctx, msg, receivers); err != nil {
		logger.Err(ctx, err, "")
		return ErrConvIndex
	}

	var (
		onlineUserDeviceSlice []string
		offlineUIDSlice       []string
	)
	for _, uid := range receivers {
		online, _ := global.Redis.SMembers(ctx, svc.RouterSrv.genRouteKey(uid)).Result()
		if len(online) > 0 {
			onlineUserDeviceSlice = append(onlineUserDeviceSlice, online...)
		} else {
			offlineUIDSlice = append(offlineUIDSlice, uid)
		}
	}
	if len(onlineUserDeviceSlice) > 0 {
		err := svc.RouterSrv.RouteChat(ctx, onlineUserDeviceSlice, msg)
		if err != nil {
			logger.Err(ctx, err, "")
			return ErrMsgRoute
		}
	}
	if len(offlineUIDSlice) > 0 {
		// todo 消息通知（离线）
		logger.Info(ctx, "消息离线通知")
	}
	return nil
}
