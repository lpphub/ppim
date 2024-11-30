package rpc

import (
	"fmt"
	"github.com/lpphub/golib/logger"
	"github.com/rpcxio/rpcx-etcd/serverplugin"
	"github.com/smallnest/rpcx/server"
	"time"
)

type RpcServer struct {
	addr     string
	etcdAddr []string
	etcdPath string
	srv      *server.Server
}

func NewRpcServer() *RpcServer {
	return &RpcServer{
		addr:     ":9090",
		etcdAddr: []string{"localhost:2379"},
		etcdPath: "/rpcx",
		srv:      server.NewServer(),
	}
}

func (s *RpcServer) registerServer() {
	_ = s.srv.RegisterName("logic", new(logicService), "")
}

func (s *RpcServer) Start() {
	setupEtcdRegisterPlugin(s)
	s.registerServer()

	logger.Log().Info().Msgf("Listening and serving RPC on %s", s.addr)
	if err := s.srv.Serve("tcp", s.addr); err != nil {
		panic(fmt.Sprintf("rpc server start failed, err:%v\n", err))
	}
}

func setupEtcdRegisterPlugin(s *RpcServer) {
	r := &serverplugin.EtcdV3RegisterPlugin{
		ServiceAddress: fmt.Sprintf("tcp@localhost%s", s.addr),
		EtcdServers:    s.etcdAddr,
		BasePath:       s.etcdPath,
		UpdateInterval: time.Minute,
		//Metrics:        metrics.NewRegistry(),
	}
	if err := r.Start(); err != nil {
		logger.Log().Err(err).Msg("failed to start etcd")
	}
	s.srv.Plugins.Add(r)
}

func (s *RpcServer) Stop() {
	_ = s.srv.Close()
}
