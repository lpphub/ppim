package net

type Server struct {
	tcpServer *TCPServer
	wsServer  *WsServer
}

func NewServer(svc *ServerContext, tcpAddr, wsAddr string) *Server {
	s := new(Server)
	if tcpAddr != "" {
		s.tcpServer = NewTCPServer(svc, tcpAddr)
	}
	if wsAddr != "" {
		s.wsServer = NewWsServer(svc, wsAddr)
	}
	return s
}

func (s *Server) Start() {
	// 启动tcp服务
	if s.tcpServer != nil {
		go func() {
			if err := s.tcpServer.Start(); err != nil {
				panic(err.Error())
			}
		}()
	}

	// 启动websocket服务
	if s.wsServer != nil {
		go func() {
			if err := s.wsServer.Start(); err != nil {
				panic(err.Error())
			}
		}()
	}
}

func (s *Server) Stop() {
	if s.tcpServer != nil {
		s.tcpServer.Stop()
	}
	if s.wsServer != nil {
		s.wsServer.Stop()
	}
}
