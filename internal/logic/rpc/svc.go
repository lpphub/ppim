package rpc

import (
	"context"
	"fmt"
	"ppim/api/logic"
	"ppim/internal/logic/global"
	"ppim/internal/logic/model"
)

type logicService struct {
	logic.UnimplementedLogicServer
}

const (
	cacheOnlineUid = "online:%s"
)

func (s *logicService) Auth(ctx context.Context, req *logic.AuthReq) (*logic.AuthResp, error) {
	user := new(model.User)
	if err := user.GetOne(ctx, req.Uid); err != nil {
		return nil, err
	}
	code := 0
	if user.DID != req.Did || user.Token != req.Token {
		code = 1 // 鉴权失败
	}
	resp := &logic.AuthResp{
		Code: uint32(code),
	}
	return resp, nil
}

// Register 登记接入层的用户连接信息
func (s *logicService) Register(ctx context.Context, req *logic.OnlineReq) (*logic.OnlineResp, error) {
	onlineVal := s.buildOnlineVal(ctx, req)
	err := global.Redis.SAdd(ctx, fmt.Sprintf(cacheOnlineUid, req.Uid), onlineVal).Err()
	if err != nil {
		return nil, err
	}
	return &logic.OnlineResp{Code: 0}, nil
}

func (s *logicService) UnRegister(ctx context.Context, req *logic.OnlineReq) (*logic.OnlineResp, error) {
	onlineVal := s.buildOnlineVal(ctx, req)
	err := global.Redis.SRem(ctx, fmt.Sprintf(cacheOnlineUid, req.Uid), onlineVal).Err()
	if err != nil {
		return nil, err
	}
	return &logic.OnlineResp{Code: 0}, nil
}

func (*logicService) buildOnlineVal(_ context.Context, req *logic.OnlineReq) string {
	return fmt.Sprintf("%s_%s_%s_%s", req.Uid, req.Did, req.Topic, req.Ip)
}
