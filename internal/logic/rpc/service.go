package rpc

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/lpphub/golib/logger"
	"ppim/api/rpctypes"
	"ppim/internal/logic/model"
	"ppim/internal/logic/srv"
	"ppim/internal/logic/types"
)

type logicService struct{}

func (s *logicService) Auth(ctx context.Context, req *rpctypes.AuthReq, resp *rpctypes.AuthResp) error {
	user := new(model.User)
	if err := user.GetOne(ctx, req.Uid); err != nil {
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
	return srv.OnSrv.Register(ctx, &ol)
}

func (s *logicService) UnRegister(ctx context.Context, req *rpctypes.RouterReq, _ *rpctypes.RouterResp) error {
	var ol types.OnlineDTO
	_ = copier.Copy(&ol, req)
	return srv.OnSrv.UnRegister(ctx, &ol)
}

func (s *logicService) SendMsg(ctx context.Context, req *rpctypes.MessageReq, _ *rpctypes.MessageResp) error {
	logger.Infof(ctx, "send msg param: %v", req)
	var msg types.MessageDTO
	_ = copier.Copy(&msg, req)
	return srv.MsgSrv.HandleMsg(ctx, &msg)
}
