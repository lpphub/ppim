package rpc

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/lpphub/golib/logger"
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
	var ol types.OnlineDTO
	_ = copier.Copy(&ol, req)
	return service.Inst().OnlineSrv.Register(ctx, &ol)
}

func (s *logicService) UnRegister(ctx context.Context, req *rpctypes.RouterReq, _ *rpctypes.RouterResp) error {
	var ol types.OnlineDTO
	_ = copier.Copy(&ol, req)
	return service.Inst().OnlineSrv.UnRegister(ctx, &ol)
}

func (s *logicService) SendMsg(ctx context.Context, req *rpctypes.MessageReq, _ *rpctypes.MessageResp) error {
	logger.Infof(ctx, "send msg param: %v", req)
	var msg types.MessageDTO
	_ = copier.Copy(&msg, req)
	return service.Inst().MsgSrv.HandleMsg(ctx, &msg)
}
