package gate

import (
	"ppim/internal/gate/global"
	"ppim/internal/gate/net"
	"ppim/internal/gate/rpc"
)

func Serve() {
	global.Init()

	if err := rpc.RegisterGrpcClient(global.Conf.Grpc.Addr); err != nil {
		panic(err.Error())
	}

	tcp := net.NewTCPServer(":5050")
	if err := tcp.Start(); err != nil {
		panic(err.Error())
	}
}
