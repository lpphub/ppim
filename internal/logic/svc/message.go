package svc

import (
	"context"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/lpphub/golib/logger"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"ppim/internal/chatlib"
	"ppim/internal/logic/global"
	"ppim/internal/logic/store"
	"ppim/internal/logic/svc/seq"
	"ppim/internal/logic/types"
	"time"
)

var (
	ErrMsgStore      = errors.New("消息持久化失败")
	ErrConvIndex     = errors.New("消息会话索引失败")
	ErrRouteDelivery = errors.New("消息路由投递失败")
)

const (
	CacheGroupMembers = "group:member:%s"
)

type MessageSrv struct {
	conv  *ConversationSrv
	route *RouteSrv
	seq   seq.Sequence
}

func newMessageSrv(conv *ConversationSrv, route *RouteSrv) *MessageSrv {
	return &MessageSrv{
		conv:  conv,
		route: route,
		seq:   seq.NewStepSequence(100, seq.WithDefaultSeqStorage(global.Redis, new(store.Conversation))),
	}
}

func (s *MessageSrv) HandleMessage(ctx context.Context, msg *types.MessageDTO) error {
	ctx = logger.WithCtx(ctx)

	// 同一会话的消息序列号是递增的
	msgSeq, err := s.seq.Next(ctx, msg.ConversationID)
	if err != nil {
		logger.Err(ctx, err, "generate msg_seq err")
		return err
	}
	msg.MsgSeq = msgSeq
	msg.UpdatedAt = time.Now().UnixMilli()

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

	// 2.获取接收者（包含自身多端）
	receivers := []string{msg.FromUID}
	if msg.ConversationType == chatlib.ConvSingle && msg.ToID != msg.FromUID {
		receivers = append(receivers, msg.ToID)
	} else if msg.ConversationType == chatlib.ConvGroup {
		members, gerr := s.getGroupMembers(ctx, msg.ToID)
		if gerr != nil {
			logger.Err(ctx, gerr, "")
		} else {
			receivers = append(receivers, members...)
		}
	}

	// 3.索引会话最新消息
	if err = s.conv.IndexUpdate(ctx, msg, receivers); err != nil {
		logger.Err(ctx, err, "conv index recent")
		return ErrConvIndex
	}

	// 4.消息投递（在线、离线）
	if err = s.route.RouteChat(ctx, msg, receivers); err != nil {
		logger.Err(ctx, err, "route chat delivery")
		return ErrRouteDelivery
	}
	return nil
}

func (s *MessageSrv) getGroupMembers(ctx context.Context, groupID string) ([]string, error) {
	members := global.Redis.SMembers(ctx, fmt.Sprintf(CacheGroupMembers, groupID)).Val()
	if len(members) > 0 {
		return members, nil
	}
	return new(store.Group).ListMembers(ctx, groupID)
}

func (s *MessageSrv) PullUpOrDown(ctx context.Context, conversationID string, startSeq, limit int64) ([]types.MessageDTO, error) {
	list, err := new(store.Message).PullUpOrDown(ctx, conversationID, startSeq, limit)
	if err != nil {
		return nil, err
	}

	voList := make([]types.MessageDTO, 0, len(list))
	for _, v := range list {
		var vo types.MessageDTO
		_ = copier.Copy(&vo, v)
		vo.SendTime = v.SendTime.UnixMilli()
		vo.CreatedAt = v.CreatedAt.UnixMilli()
		vo.UpdatedAt = v.UpdatedAt.UnixMilli()
		voList = append(voList, vo)
	}
	return voList, nil
}

func (s *MessageSrv) Revoke(ctx context.Context, msgID, conversationID string) error {
	if err := new(store.Message).Revoke(ctx, msgID); err != nil {
		return err
	}
	// 如果撤回的消息是会话最新消息，则更新会话最新消息
	msg, err := s.conv.cacheQueryLastMsg(ctx, conversationID)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil
		}
		return err
	}
	if msg.MsgID == msgID {
		msg.Revoked = 1
		msg.UpdatedAt = time.Now().UnixMilli()
		s.conv.cacheStoreLastMsg(ctx, msg)
	}
	// todo 同步事件，会话未读数、最新时间是否要变？
	return nil
}

func (s *MessageSrv) Delete(ctx context.Context, msgID, conversationID string) error {
	if err := new(store.Message).Delete(ctx, msgID); err != nil {
		return err
	}
	// 如果删除的消息是会话最新消息，则更新会话最新消息
	msg, err := s.conv.cacheQueryLastMsg(ctx, conversationID)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil
		}
		return err
	}
	if msg.MsgID == msgID {
		msg.Deleted = 1
		msg.UpdatedAt = time.Now().UnixMilli()
		s.conv.cacheStoreLastMsg(ctx, msg)
	}
	// todo 同步事件，会话未读数、最新时间是否要变？
	return nil
}
