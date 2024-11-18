package rpc

import (
	"context"
	"google.golang.org/grpc"
	"ppim/api/logic"
)

type AuthSrv struct {
	conn *grpc.ClientConn //global conn
}

func (srv *AuthSrv) Auth(ctx context.Context, uid, did, token string) (bool, error) {
	resp, err := logic.NewAuthClient(srv.conn).Auth(ctx, &logic.AuthReq{
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
