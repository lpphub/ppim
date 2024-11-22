package rpc

import (
	"context"
	"ppim/api/logic"
)

type OnlineSrv struct {
}

func (srv *OnlineSrv) Register(ctx context.Context, uid, did, topic, ip string) (bool, error) {
	resp, err := logic.NewOnlineClient(caller.conn).Register(ctx, &logic.OnlineReq{
		Uid:   uid,
		Did:   did,
		Topic: topic,
		Ip:    ip,
	})
	if err != nil {
		return false, err
	}
	return resp.Ok, nil
}

func (srv *OnlineSrv) UnRegister(ctx context.Context, uid, did, topic, ip string) (bool, error) {
	resp, err := logic.NewOnlineClient(caller.conn).Unregister(ctx, &logic.OnlineReq{
		Uid:   uid,
		Did:   did,
		Topic: topic,
		Ip:    ip,
	})
	if err != nil {
		return false, err
	}
	return resp.Ok, nil
}
