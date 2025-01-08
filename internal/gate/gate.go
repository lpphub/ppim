package gate

import (
	"os"
	"os/signal"
	"ppim/internal/gate/global"
	"ppim/internal/gate/mq"
	"ppim/internal/gate/net"
	"ppim/internal/gate/rpc"
	"syscall"
)

func Serve() {
	global.Init()

	// 注册rpc客户端
	if err := rpc.RegisterRpcClient(global.Conf.Server.RpcRegistry); err != nil {
		panic(err.Error())
	}

	// 初始化服务上下文
	svc := net.InitServerContext()

	// 订阅消息
	if err := mq.RegisterSubscriber(svc); err != nil {
		panic(err.Error())
	}

	// 启动tcp服务
	go func() {
		tcp := net.NewTCPServer(svc, global.Conf.Server.Tcp)
		if err := tcp.Start(); err != nil {
			panic(err.Error())
		}
		defer tcp.Stop()
	}()

	go func() {
		ws := net.NewWsServer(svc, global.Conf.Server.Ws)
		if err := ws.Start(); err != nil {
			panic(err.Error())
		}
		defer ws.Stop()
	}()

	sign := make(chan os.Signal, 1)
	signal.Notify(sign, syscall.SIGINT, syscall.SIGTERM)

	// 等待信号
	<-sign
}
