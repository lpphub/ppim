package gate

import (
	"os"
	"os/signal"
	"ppim/internal/gate/global"
	"ppim/internal/gate/net"
	"ppim/internal/gate/rpc"
	"ppim/internal/gate/sub"
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
	if err := sub.SubscribeDelivery(svc); err != nil {
		panic(err.Error())
	}

	server := net.NewServer(svc, global.Conf.Server.Tcp, global.Conf.Server.Ws)
	server.Start()
	defer server.Stop()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	// 等待信号
	<-sigChan
}
