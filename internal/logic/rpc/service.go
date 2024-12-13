package rpc

import (
	"context"
	"github.com/jinzhu/copier"
	"ppim/api/rpctypes"
	"ppim/internal/logic/service"
	"ppim/internal/logic/store"
	"ppim/internal/logic/types"
)

type logicService struct{}

func (s *logicService) Auth(ctx context.Context, req *rpctypes.AuthReq, resp *rpctypes.AuthResp) error {
	user, err := new(store.User).GetOne(ctx, req.Uid)
	if err != nil {
		return err
	}
	code := 0
	if user.DID != req.Did || user.Token != req.Token {
		code = 1 // 鉴权失败
	}
	*resp = rpctypes.AuthResp{
		Code: code,
	}
	return nil
}

func (s *logicService) Register(ctx context.Context, req *rpctypes.RouterReq, _ *rpctypes.RouterResp) error {
	var ol types.RouteDTO
	_ = copier.Copy(&ol, req)
	return service.Hint().Route.Online(ctx, &ol)
}

func (s *logicService) UnRegister(ctx context.Context, req *rpctypes.RouterReq, _ *rpctypes.RouterResp) error {
	var ol types.RouteDTO
	_ = copier.Copy(&ol, req)
	return service.Hint().Route.Offline(ctx, &ol)
}

func (s *logicService) SendMsg(ctx context.Context, req *rpctypes.MessageReq, _ *rpctypes.MessageResp) error {
	var msg types.MessageDTO
	_ = copier.Copy(&msg, req)
	return service.Hint().Msg.HandleMsg(ctx, &msg)
}
