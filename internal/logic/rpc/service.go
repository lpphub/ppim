package rpc

import (
	"context"
	"github.com/jinzhu/copier"
	"ppim/internal/chatlib"
	"ppim/internal/logic/store"
	"ppim/internal/logic/svc"
	"ppim/internal/logic/types"
)

type logicService struct{}

func (s *logicService) Auth(ctx context.Context, req *chatlib.AuthReq, _ *chatlib.AuthResp) error {
	// todo 先从缓存验证
	_, err := new(store.User).GetByToken(ctx, req.Uid, req.Token)
	if err != nil {
		return err
	}
	return nil
}

func (s *logicService) Register(ctx context.Context, req *chatlib.RouterReq, _ *chatlib.RouterResp) error {
	var ol types.RouteDTO
	_ = copier.Copy(&ol, req)
	return svc.Hints().Route.Online(ctx, &ol)
}

func (s *logicService) UnRegister(ctx context.Context, req *chatlib.RouterReq, _ *chatlib.RouterResp) error {
	var ol types.RouteDTO
	_ = copier.Copy(&ol, req)
	return svc.Hints().Route.Offline(ctx, &ol)
}

func (s *logicService) SendMsg(ctx context.Context, req *chatlib.MessageReq, resp *chatlib.MessageResp) error {
	var msg types.MessageDTO
	_ = copier.Copy(&msg, req)

	err := svc.Hints().Msg.HandleMessage(ctx, &msg)
	if err != nil {
		return err
	}
	resp.MsgSeq = msg.MsgSeq
	return nil
}
