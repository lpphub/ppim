package net

import (
	"fmt"
)

type TCPServer struct {
	Addr      string
	Connected int32

	connManager *ClientManager
	engine      *EventEngine
}

type ServerContext struct {
	*TCPServer
}

func NewTCPServer(addr string) *TCPServer {
	server := &TCPServer{
		Addr:        fmt.Sprintf("tcp://%s", addr),
		connManager: newClientManager(),
	}

	context := &ServerContext{
		TCPServer: server,
	}

	server.engine = newEventEngine(context)
	return server
}

func (s *TCPServer) Start() error {
	return s.engine.start()
}
