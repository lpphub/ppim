package rpc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcCaller struct {
	AuthClient *AuthSrv
}

var (
	Caller *GrpcCaller
)

func RegisterGrpcClient(addr string) error {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	Caller = &GrpcCaller{
		AuthClient: &AuthSrv{conn: conn},
	}
	return nil
}
