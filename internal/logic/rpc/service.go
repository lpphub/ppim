package rpc

import (
	"context"
	"github.com/jinzhu/copier"
	"ppim/internal/chatlib"
	"ppim/internal/logic/service"
	"ppim/internal/logic/store"
	"ppim/internal/logic/types"
)

type logicService struct{}

func (s *logicService) Auth(ctx context.Context, req *chatlib.AuthReq, resp *chatlib.AuthResp) error {
	user, err := new(store.User).GetOne(ctx, req.Uid)
	if err != nil {
		return err
	}
	code := 0
	if user.Token != req.Token {
		code = 1 // 鉴权失败
	}
	*resp = chatlib.AuthResp{
		Code: code,
	}
	return nil
}

func (s *logicService) Register(ctx context.Context, req *chatlib.RouterReq, _ *chatlib.RouterResp) error {
	var ol types.RouteDTO
	_ = copier.Copy(&ol, req)
	return service.Hints().Route.Online(ctx, &ol)
}

func (s *logicService) UnRegister(ctx context.Context, req *chatlib.RouterReq, _ *chatlib.RouterResp) error {
	var ol types.RouteDTO
	_ = copier.Copy(&ol, req)
	return service.Hints().Route.Offline(ctx, &ol)
}

func (s *logicService) SendMsg(ctx context.Context, req *chatlib.MessageReq, resp *chatlib.MessageResp) error {
	var msg types.MessageDTO
	_ = copier.Copy(&msg, req)

	err := service.Hints().Msg.HandleMsg(ctx, &msg)
	if err != nil {
		return err
	}
	resp.MsgSeq = msg.MsgSeq
	return nil
}
