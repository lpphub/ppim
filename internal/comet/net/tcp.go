package net

import (
	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/logging"
	"sync/atomic"
)

type TCPServer struct {
	gnet.BuiltinEventEngine
	engine gnet.Engine

	Addr string

	connected int32
	Conns     map[uint64]*gnet.Conn
}

func NewTCPServer(addr string) *TCPServer {
	return &TCPServer{
		Addr:  addr,
		Conns: make(map[uint64]*gnet.Conn),
	}
}

func (s *TCPServer) Start() error {
	return gnet.Run(s, s.Addr, gnet.WithMulticore(true))
}

func (s *TCPServer) OnBoot(eng gnet.Engine) gnet.Action {
	s.engine = eng
	logging.Infof("Listening and serving TCP on %s\n", s.Addr)
	return gnet.None
}

func (s *TCPServer) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	// 读取连接认证信息
	buf, _ := c.Next(-1)
	logging.Infof("OnOpen: %s\n", string(buf))

	atomic.AddInt32(&s.connected, 1)

	return
}

func (s *TCPServer) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	if err != nil {
		logging.Infof("error occurred on connection=%s, %v\n", c.RemoteAddr().String(), err)
	}
	atomic.AddInt32(&s.connected, -1)

	// todo 连接关闭

	return
}

func (s *TCPServer) OnTraffic(c gnet.Conn) gnet.Action {
	buf, _ := c.Next(-1)

	c.Write(buf)
	return gnet.None
}
