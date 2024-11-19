package comet

import (
	"ppim/internal/comet/global"
	"ppim/internal/comet/net"
	"ppim/internal/comet/rpc"
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
