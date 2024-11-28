package rpc

import (
	"context"
	"github.com/lpphub/golib/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"ppim/api/logic"
	"sync"
)

type GrpcCaller struct {
	logic logic.LogicClient
}

var (
	caller *GrpcCaller
	once   sync.Once
)

func RegisterGrpcClient(addr string) (err error) {
	once.Do(func() {
		_conn, cerr := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if cerr != nil {
			err = cerr
			return
		}

		caller = &GrpcCaller{
			logic: logic.NewLogicClient(_conn),
		}
	})
	return
}

func Caller() *GrpcCaller {
	return caller
}

func Context() context.Context {
	return logger.WithCtx(context.Background())
}

func (c *GrpcCaller) Auth(ctx context.Context, uid, did, token string) (bool, error) {
	resp, err := c.logic.Auth(ctx, &logic.AuthReq{
		Uid:   uid,
		Did:   did,
		Token: token,
	})
	if err != nil {
		logger.Err(ctx, err, "rpc - auth error")
		return false, err
	}
	return resp.Code == 0, nil
}

func (c *GrpcCaller) Register(ctx context.Context, uid, did string) error {
	// todo: 获取当前节点对应的topic及ip
	_, err := c.logic.Register(ctx, &logic.OnlineReq{
		Uid:   uid,
		Did:   did,
		Topic: "topic",
		Ip:    "ip",
	})
	if err != nil {
		logger.Err(ctx, err, "")
		return err
	}
	return nil
}

func (c *GrpcCaller) UnRegister(ctx context.Context, uid, did string) error {
	_, err := c.logic.Unregister(ctx, &logic.OnlineReq{
		Uid:   uid,
		Did:   did,
		Topic: "topic",
		Ip:    "ip",
	})
	if err != nil {
		logger.Err(ctx, err, "")
		return err
	}
	return nil
}
