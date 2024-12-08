package rpc

import (
	"context"
	"github.com/lpphub/golib/logger"
	etcdclient "github.com/rpcxio/rpcx-etcd/client"
	"github.com/smallnest/rpcx/client"
	"ppim/api/rpctypes"
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
	etcdPath    = "/rpcx"
	serviceName = "logic"

	methodAuth       = "Auth"
	methodRegister   = "Register"
	methodUnRegister = "UnRegister"
	methodSendMsg    = "SendMsg"
)

func RegisterRpcClient(addr string) (err error) {
	once.Do(func() {
		discovery, derr := etcdclient.NewEtcdV3Discovery(etcdPath, serviceName, strings.Split(addr, ","), true, nil)
		if derr != nil {
			err = derr
			return
		}
		logic := client.NewXClient(serviceName, client.Failtry, client.RandomSelect, discovery, client.DefaultOption)

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
	req := &rpctypes.AuthReq{
		Uid:   uid,
		Did:   did,
		Token: token,
	}
	resp := &rpctypes.AuthResp{}

	err := c.logic.Call(ctx, methodAuth, req, resp)
	if err != nil {
		logger.Err(ctx, err, "rpc - auth error")
		return false, err
	}
	return resp.Code == 0, nil
}

func (c *RpcCaller) Register(ctx context.Context, uid, did string) error {
	// todo: 获取当前节点对应的topic及ip
	req := &rpctypes.RouterReq{
		Uid: uid,
		Did: did,
	}
	resp := &rpctypes.RouterReq{}

	err := c.logic.Call(ctx, methodRegister, req, resp)
	if err != nil {
		logger.Err(ctx, err, "")
		return err
	}
	return nil
}

func (c *RpcCaller) UnRegister(ctx context.Context, uid, did string) error {
	req := &rpctypes.RouterReq{
		Uid: uid,
		Did: did,
	}
	resp := &rpctypes.RouterReq{}

	err := c.logic.Call(ctx, methodUnRegister, req, resp)
	if err != nil {
		logger.Err(ctx, err, "")
		return err
	}
	return nil
}

func (c *RpcCaller) SendMsg(ctx context.Context, msg *rpctypes.MessageReq) (*rpctypes.MessageResp, error) {
	req := msg
	resp := &rpctypes.MessageResp{}

	err := c.logic.Call(ctx, methodSendMsg, req, resp)
	if err != nil {
		logger.Err(ctx, err, "")
		return nil, err
	}
	return resp, nil
}
