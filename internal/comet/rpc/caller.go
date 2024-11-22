package rpc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcCaller struct {
	AuthClient *AuthSrv
}

var (
	caller *GrpcCaller
)

func RegisterGrpcClient(addr string) error {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	caller = &GrpcCaller{
		AuthClient: &AuthSrv{conn: conn},
	}
	return nil
}

func Caller() *GrpcCaller {
	return caller
}
