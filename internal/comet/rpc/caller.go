package rpc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcCaller struct {
	conn *grpc.ClientConn //global conn

	AuthSrv   *AuthSrv
	OnlineSrv *OnlineSrv
}

var (
	caller *GrpcCaller
)

func RegisterGrpcClient(addr string) error {
	_conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	caller = &GrpcCaller{
		conn:      _conn,
		AuthSrv:   &AuthSrv{},
		OnlineSrv: &OnlineSrv{},
	}
	return nil
}

func Caller() *GrpcCaller {
	return caller
}
