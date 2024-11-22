package rpc

import (
	"context"
	"fmt"
	"ppim/api/logic"
	"ppim/internal/logic/global"
)

type onlineService struct {
	logic.UnimplementedOnlineServer
}

const (
	cacheOnlineUid = "online:%s"
)

// Register 登记接入层的用户连接信息
func (ol *onlineService) Register(ctx context.Context, req *logic.OnlineReq) (*logic.OnlineResp, error) {
	onlineVal := ol.buildOnlineVal(ctx, req)
	err := global.Redis.SAdd(ctx, fmt.Sprintf(cacheOnlineUid, req.Uid), onlineVal).Err()
	if err != nil {
		return nil, err
	}
	return &logic.OnlineResp{Ok: true}, nil
}

func (ol *onlineService) UnRegister(ctx context.Context, req *logic.OnlineReq) (*logic.OnlineResp, error) {
	onlineVal := ol.buildOnlineVal(ctx, req)
	err := global.Redis.SRem(ctx, fmt.Sprintf(cacheOnlineUid, req.Uid), onlineVal).Err()
	if err != nil {
		return nil, err
	}
	return &logic.OnlineResp{Ok: true}, nil
}

func (*onlineService) buildOnlineVal(_ context.Context, req *logic.OnlineReq) string {
	return fmt.Sprintf("%s_%s_%s_%s", req.Uid, req.Did, req.Topic, req.Ip)
}
