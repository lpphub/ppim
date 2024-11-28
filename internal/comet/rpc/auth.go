package rpc

import (
	"context"
	"github.com/lpphub/golib/logger"
	"ppim/api/logic"
)

type AuthSrv struct {
}

func (srv *AuthSrv) Auth(ctx context.Context, uid, did, token string) (bool, error) {
	resp, err := logic.NewAuthClient(caller.conn).Auth(ctx, &logic.AuthReq{
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
