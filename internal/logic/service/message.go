package service

import (
	"context"
	"fmt"
	"github.com/lpphub/golib/logger"
	"github.com/pkg/errors"
	"ppim/internal/chatlib"
	"ppim/internal/logic/service/seq"
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
	seq   seq.Sequence
}

func newMessageSrv(conv *ConversationSrv, route *RouteSrv, seq seq.Sequence) *MessageSrv {
	return &MessageSrv{
		conv:  conv,
		route: route,
		seq:   seq,
	}
}

func (s *MessageSrv) HandleMsg(ctx context.Context, msg *types.MessageDTO) error {
	ctx = logger.WithCtx(ctx)

	msgSeq, err := s.seq.Next(ctx, msg.ConversationID)
	if err != nil {
		logger.Err(ctx, err, "generate msg_seq err")
		return err
	}
	msg.MsgSeq = msgSeq

	// todo 优化：异步处理

	// 1.消息持久化
	mm := &store.Message{
		MsgNo:            msg.MsgNo,
		MsgID:            msg.MsgID,
		MsgSeq:           msg.MsgSeq,
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
	if err = mm.Insert(ctx); err != nil {
		logger.Err(ctx, err, "")
		return ErrMsgStore
	}

	// 2.获取接收者（包含自身）
	receivers := []string{msg.FromUID}
	if msg.ConversationType == chatlib.ConvSingle && msg.ToID != msg.FromUID {
		receivers = append(receivers, msg.ToID)
	} else if msg.ConversationType == chatlib.ConvGroup {
		members, serr := new(store.Group).ListMembers(ctx, msg.ToID)
		if serr != nil {
			logger.Err(ctx, serr, "")
		} else {
			receivers = append(receivers, members...)
		}
	}

	// 3.索引会话最新消息
	if err = s.conv.IndexRecent(ctx, msg, receivers); err != nil {
		logger.Err(ctx, err, "")
		return ErrConvIndex
	}

	// 4.在线投递
	var (
		onlineSlice  []string //在线用户路由
		offlineSlice []string //离线用户UID
	)
	receiverChunks := util.SplitSlice(receivers, 500)
	for _, chunks := range receiverChunks {
		cmds, berr := s.route.BatchGetOnline(ctx, chunks)
		if berr != nil {
			logger.Err(ctx, berr, fmt.Sprintf("batch get online error: %v", chunks))
			continue
		}

		for i, uid := range chunks {
			online, _ := cmds[i].Result()
			if len(online) > 0 {
				for did, topic := range online {
					if did == msg.FromDID && uid == msg.FromUID { // 排除发送者同一设备，而不同设备时则接收消息
						continue
					}
					onlineSlice = append(onlineSlice, fmt.Sprintf("%s#%s", uid, topic))
				}
			} else {
				offlineSlice = append(offlineSlice, uid)
			}
		}
	}

	if len(onlineSlice) > 0 {
		err = s.route.RouteDelivery(ctx, util.RemoveDup(onlineSlice), msg)
		if err != nil {
			logger.Err(ctx, err, "")
			return ErrMsgRoute
		}
	}

	if len(offlineSlice) > 0 {
		// todo 消息离线通知
		logger.Warnf(ctx, "offline push: %v", offlineSlice)
	}
	return nil
}
