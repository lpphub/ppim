package rpc

import (
	"context"
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
		// TODO: log
		return false, err
	}
	return resp.Ok, nil
}
