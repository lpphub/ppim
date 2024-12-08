package gate

import (
	"ppim/internal/gate/global"
	"ppim/internal/gate/mq"
	"ppim/internal/gate/net"
	"ppim/internal/gate/rpc"
)

func Serve() {
	global.Init()

	if err := rpc.RegisterRpcClient(global.Conf.Server.Etcd); err != nil {
		panic(err.Error())
	}

	mq.RegisterSubscriber()

	tcp := net.NewTCPServer(global.Conf.Server.Tcp)
	if err := tcp.Start(); err != nil {
		panic(err.Error())
	}
	defer tcp.Stop()
}
