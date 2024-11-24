package rpc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"sync"
)

type GrpcCaller struct {
	conn *grpc.ClientConn //global conn

	AuthSrv   *AuthSrv
	OnlineSrv *OnlineSrv
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
			conn:      _conn,
			AuthSrv:   &AuthSrv{},
			OnlineSrv: &OnlineSrv{},
		}
	})
	return
}

func Caller() *GrpcCaller {
	return caller
}
