package rpc

import (
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"ppim/api/logic"
)

type GrpcServer struct {
	srv *grpc.Server
}

func NewGrpcServer() *GrpcServer {
	return &GrpcServer{
		srv: grpc.NewServer(),
	}
}

func (s *GrpcServer) registerServer() {
	logic.RegisterPingServer(s.srv, &pingServer{})
}

func (s *GrpcServer) Start() {
	s.registerServer()

	lis, err := net.Listen("tcp", ":9090")
	if err != nil {
		panic(fmt.Sprintf("grpc server listen err:%v\n", err))
	}

	log.Printf("Listening and serving GRPC on %s\n", lis.Addr().String())
	if err = s.srv.Serve(lis); err != nil {
		panic(fmt.Sprintf("grpc server start failed, err:%v\n", err))
	}
}

func (s *GrpcServer) Stop() {
	s.srv.GracefulStop()
}
