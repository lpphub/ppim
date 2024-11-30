package rpc

import (
	"context"
	"fmt"
	"ppim/api/rpctypes"
	"ppim/internal/logic/global"
	"ppim/internal/logic/model"
)

type logicService struct{}

const (
	cacheRouteUid = "router:%s"
)

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
	onlineVal := buildRouterVal(ctx, req)
	err := global.Redis.SAdd(ctx, fmt.Sprintf(cacheRouteUid, req.Uid), onlineVal).Err()
	if err != nil {
		return err
	}
	return nil
}

func (s *logicService) UnRegister(ctx context.Context, req *rpctypes.RouterReq, _ *rpctypes.RouterResp) error {
	onlineVal := buildRouterVal(ctx, req)
	err := global.Redis.SRem(ctx, fmt.Sprintf(cacheRouteUid, req.Uid), onlineVal).Err()
	if err != nil {
		return err
	}
	return nil
}

func buildRouterVal(_ context.Context, req *rpctypes.RouterReq) string {
	return fmt.Sprintf("%s_%s_%s_%s", req.Uid, req.Did, req.Topic, req.Ip)
}
