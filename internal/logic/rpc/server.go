package rpc

import (
	"fmt"
	"github.com/lpphub/golib/logger"
	"github.com/rpcxio/rpcx-etcd/serverplugin"
	"github.com/smallnest/rpcx/server"
	"strings"
	"time"
)

type RpcServer struct {
	serviceAddr  string
	registryAddr []string
	basePath     string
	srv          *server.Server
}

func NewRpcServer(srvAddr, registry string) *RpcServer {
	return &RpcServer{
		serviceAddr:  srvAddr,
		registryAddr: strings.Split(registry, ","),
		basePath:     "/rpcx",
		srv:          server.NewServer(),
	}
}

func (s *RpcServer) registerServer() {
	_ = s.srv.RegisterName("logic", new(logicService), "")
}

func (s *RpcServer) Start() {
	setupRegistryPlugin(s)
	s.registerServer()

	logger.Log().Info().Msgf("Listening and serving RPC on %s", s.serviceAddr)
	if err := s.srv.Serve("tcp", s.serviceAddr); err != nil {
		panic(fmt.Sprintf("rpc server start failed, err:%v\n", err))
	}
}

func setupRegistryPlugin(s *RpcServer) {
	r := &serverplugin.EtcdV3RegisterPlugin{
		ServiceAddress: fmt.Sprintf("tcp@%s", s.serviceAddr),
		EtcdServers:    s.registryAddr,
		BasePath:       s.basePath,
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
