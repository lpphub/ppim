package net

import "fmt"

type TCPServer struct {
	context *ServerContext
	engine  *EventEngine
}

type ServerContext struct {
	Addr   string
	online int32

	ConnManager *ClientManager
}

func NewTCPServer(addr string) *TCPServer {
	context := &ServerContext{
		Addr:        fmt.Sprintf("tcp://%s", addr),
		ConnManager: newClientManager(),
	}

	server := &TCPServer{
		context: context,
		engine:  newEventEngine(context),
	}
	return server
}

func (s *TCPServer) Start() error {
	return s.engine.start()
}

func (s *TCPServer) GetSvc() *ServerContext {
	return s.context
}
