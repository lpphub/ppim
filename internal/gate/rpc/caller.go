package rpc

import (
	"context"
	"fmt"
	"github.com/lpphub/golib/logger"
	etcdclient "github.com/rpcxio/rpcx-etcd/client"
	"github.com/smallnest/rpcx/client"
	"ppim/internal/chatlib"
	"ppim/internal/gate/global"
	"strings"
	"sync"
)

type RpcCaller struct {
	logic client.XClient
}

var (
	caller *RpcCaller
	once   sync.Once
)

const (
	basePath    = "/rpcx"
	serviceName = "logic"

	methodAuth       = "Auth"
	methodRegister   = "Register"
	methodUnRegister = "UnRegister"
	methodSendMsg    = "SendMsg"
)

func RegisterRpcClient(registryAddr string) (err error) {
	once.Do(func() {
		discovery, cerr := etcdclient.NewEtcdV3Discovery(basePath, serviceName, strings.Split(registryAddr, ","), true, nil)
		if cerr != nil {
			err = fmt.Errorf("register rpc discovery error: %v", cerr)
			return
		}
		logic := client.NewXClient(serviceName, client.Failtry, client.SelectByUser, discovery, client.DefaultOption)
		logic.SetSelector(&customSelector{})

		caller = &RpcCaller{
			logic: logic,
		}
	})
	return
}

func Caller() *RpcCaller {
	return caller
}

func Context() context.Context {
	return logger.WithCtx(context.Background())
}

func (c *RpcCaller) Auth(ctx context.Context, uid, did, token string) (bool, error) {
	req := &chatlib.AuthReq{
		Uid:   uid,
		Did:   did,
		Token: token,
	}
	resp := &chatlib.AuthResp{}

	err := c.logic.Call(ctx, methodAuth, req, resp)
	if err != nil {
		logger.Err(ctx, err, "rpc - auth error")
		return false, err
	}
	return resp.Code == 0, nil
}

// Register 将当前连接对应的topic注册到logic route
func (c *RpcCaller) Register(ctx context.Context, uid, did string) error {
	req := &chatlib.RouterReq{
		Uid:   uid,
		Did:   did,
		Topic: global.Conf.Kafka.Topic, // 将当前连接
	}
	err := c.logic.Call(ctx, methodRegister, req, &chatlib.RouterReq{})
	if err != nil {
		logger.Err(ctx, err, "")
		return err
	}
	return nil
}

func (c *RpcCaller) UnRegister(ctx context.Context, uid, did string) error {
	req := &chatlib.RouterReq{
		Uid:   uid,
		Did:   did,
		Topic: global.Conf.Kafka.Topic,
	}
	err := c.logic.Call(ctx, methodUnRegister, req, &chatlib.RouterReq{})
	if err != nil {
		logger.Err(ctx, err, "")
		return err
	}
	return nil
}

func (c *RpcCaller) SendMsg(ctx context.Context, msg *chatlib.MessageReq) (*chatlib.MessageResp, error) {
	req := msg
	resp := &chatlib.MessageResp{}

	err := c.logic.Call(ctx, methodSendMsg, req, resp)
	if err != nil {
		logger.Err(ctx, err, "")
		return nil, err
	}
	return resp, nil
}
