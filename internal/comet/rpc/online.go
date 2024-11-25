package rpc

import (
	"context"
	"github.com/lpphub/golib/logger"
	"ppim/api/logic"
)

type OnlineSrv struct {
}

func (srv *OnlineSrv) Register(ctx context.Context, uid, did, topic, ip string) error {
	_, err := logic.NewOnlineClient(caller.conn).Register(ctx, &logic.OnlineReq{
		Uid:   uid,
		Did:   did,
		Topic: topic,
		Ip:    ip,
	})
	if err != nil {
		logger.Err(ctx, err, "")
		return err
	}
	return nil
}

func (srv *OnlineSrv) UnRegister(ctx context.Context, uid, did, topic, ip string) error {
	_, err := logic.NewOnlineClient(caller.conn).Unregister(ctx, &logic.OnlineReq{
		Uid:   uid,
		Did:   did,
		Topic: topic,
		Ip:    ip,
	})
	if err != nil {
		logger.Err(ctx, err, "")
		return err
	}
	return nil
}
