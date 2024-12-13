package net

type WsServer struct {
	svc    *ServerContext
	engine *EventEngine
}

func NewWsServer(svc *ServerContext, addr string) *WsServer {
	engOpt := EngineOption{
		Network: _ws,
		Addr:    addr,
	}

	server := &WsServer{
		svc:    svc,
		engine: NewEventEngine(svc, engOpt),
	}
	return server
}

func (s *WsServer) Start() error {
	return s.engine.start()
}

func (s *WsServer) Stop() {
	_ = s.engine.stop()
}
