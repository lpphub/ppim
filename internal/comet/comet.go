package comet

import (
	"ppim/internal/comet/net"
)

func Serve() {
	tcp := net.NewTCPServer(":5050")
	if err := tcp.Start(); err != nil {
		panic(err.Error())
	}
}
