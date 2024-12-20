package gate

import (
	"ppim/internal/gate/global"
	"ppim/internal/gate/mq"
	"ppim/internal/gate/net"
	"ppim/internal/gate/rpc"
	"ppim/internal/gate/task"
)

func Serve() {
	global.Init()

	if err := rpc.RegisterRpcClient(global.Conf.Server.RpcRegistry); err != nil {
		panic(err.Error())
	}

	svc := net.InitServerContext()

	retry := task.NewRetryDelivery(svc)
	retry.Start()

	mq.RegisterSubscriber(svc, retry)

	go func() {
		tcp := net.NewTCPServer(svc, global.Conf.Server.Tcp)
		if err := tcp.Start(); err != nil {
			panic(err.Error())
		}
		defer tcp.Stop()
	}()

	ws := net.NewWsServer(svc, global.Conf.Server.Ws)
	if err := ws.Start(); err != nil {
		panic(err.Error())
	}
	defer ws.Stop()
}
