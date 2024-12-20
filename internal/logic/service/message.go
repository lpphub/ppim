package service

import (
	"context"
	"fmt"
	"github.com/lpphub/golib/logger"
	"github.com/pkg/errors"
	"ppim/internal/chatlib"
	"ppim/internal/logic/global"
	"ppim/internal/logic/store"
	"ppim/internal/logic/types"
	"ppim/pkg/util"
	"time"
)

var (
	ErrMsgStore  = errors.New("消息持久化失败")
	ErrConvIndex = errors.New("消息索引更新失败")
	ErrMsgRoute  = errors.New("消息路由失败")
)

type MessageSrv struct {
	conv  *ConversationSrv
	route *RouteSrv
}

func newMessageSrv(conv *ConversationSrv, route *RouteSrv) *MessageSrv {
	return &MessageSrv{
		conv:  conv,
		route: route,
	}
}

func (s *MessageSrv) HandleMsg(ctx context.Context, msg *types.MessageDTO) error {
	// todo 优化：异步处理
	ctx = logger.WithCtx(ctx)

	// 1.消息持久化
	mm := &store.Message{
		MsgID:            msg.MsgID,
		MsgSeq:           msg.MsgSeq,
		MsgNo:            msg.MsgNo,
		MsgType:          msg.MsgType,
		Content:          msg.Content,
		ConversationID:   msg.ConversationID,
		ConversationType: msg.ConversationType,
		FromUID:          msg.FromUID,
		ToID:             msg.ToID,
		SendTime:         time.UnixMilli(msg.SendTime),
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	if err := mm.Insert(ctx); err != nil {
		logger.Err(ctx, err, "")
		return ErrMsgStore
	}

	// 2.获取接收者
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
		logger.Warn(ctx, "msg receivers is empty")
		return nil
	}

	// 3.索引会话最新消息
	if err := s.conv.IndexRecent(ctx, msg, receivers); err != nil {
		logger.Err(ctx, err, "")
		return ErrConvIndex
	}

	// 4.在线投递
	var (
		onlineSlice  []string //在线用户路由
		offlineSlice []string //离线用户UID
	)
	for _, uid := range receivers {
		online, _ := global.Redis.HGetAll(ctx, s.route.genRouteKey(uid)).Result()
		if len(online) > 0 {
			for _, topic := range online {
				onlineSlice = append(onlineSlice, fmt.Sprintf("%s#%s", uid, topic))
			}
		} else {
			offlineSlice = append(offlineSlice, uid)
		}
	}
	// 发送者的其他在线设备也接收消息
	selfOnline, _ := global.Redis.HGetAll(ctx, s.route.genRouteKey(msg.FromUID)).Result()
	if len(selfOnline) > 1 {
		for did, topic := range selfOnline {
			if msg.FromDID != did {
				onlineSlice = append(onlineSlice, fmt.Sprintf("%s#%s", msg.FromUID, topic))
			}
		}
	}

	if len(onlineSlice) > 0 {
		err := s.route.RouteDelivery(ctx, util.RemoveDup(onlineSlice), msg)
		if err != nil {
			logger.Err(ctx, err, "")
			return ErrMsgRoute
		}
	}

	if len(offlineSlice) > 0 {
		// todo 消息离线通知
	}
	return nil
}
