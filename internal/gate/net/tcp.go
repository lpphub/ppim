package net

type TCPServer struct {
	svc    *ServerContext
	engine *EventEngine
}

func NewTCPServer(ctx *ServerContext, addr string) *TCPServer {
	engOpt := EngineOption{
		Protocol: _tcp,
		Addr:     addr,
	}

	server := &TCPServer{
		svc:    ctx,
		engine: NewEventEngine(ctx, engOpt),
	}
	return server
}

func (s *TCPServer) Start() error {
	return s.engine.start()
}

func (s *TCPServer) Stop() {
	_ = s.engine.stop()
}

func (s *TCPServer) GetSvc() *ServerContext {
	return s.svc
}
