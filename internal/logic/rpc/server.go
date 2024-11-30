package rpc

import (
	"fmt"
	"github.com/lpphub/golib/logger"
	"github.com/smallnest/rpcx/server"
)

type RpcServer struct {
	addr string
	srv  *server.Server
}

func NewGrpcServer() *RpcServer {
	return &RpcServer{
		addr: ":9090",
		srv:  server.NewServer(),
	}
}

func (s *RpcServer) registerServer() {
	_ = s.srv.RegisterName("logic", new(logicService), "")
}

func (s *RpcServer) Start() {
	s.registerServer()

	logger.Log().Info().Msgf("Listening and serving RPC on %s", s.addr)
	if err := s.srv.Serve("tcp", s.addr); err != nil {
		panic(fmt.Sprintf("rpc server start failed, err:%v\n", err))
	}
}

func (s *RpcServer) Stop() {
	_ = s.srv.Close()
}
