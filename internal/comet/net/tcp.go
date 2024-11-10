package net

import (
	"errors"
	"fmt"
	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/logging"
	"github.com/panjf2000/gnet/v2/pkg/pool/goroutine"
	"google.golang.org/protobuf/proto"
	"ppim/api/message_pb"
	"ppim/internal/comet/net/codec"
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
		Addr:  fmt.Sprintf("tcp://%s", addr),
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
	//c.SetContext(new(codec.FixedHeadCodec))
	c.SetContext(new(codec.ProtobufCodec))

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
	codecIns := c.Context().(*codec.ProtobufCodec)
	buf, err := codecIns.Decode(c)
	if errors.Is(err, codec.ErrInvalidMagic) { //非法数据包关闭连接，不完整的包继续处理
		return gnet.Close
	}

	_ = goroutine.Default().Submit(func() {
		// todo 异步处理业务

		var msg message_pb.Message
		_ = proto.Unmarshal(buf, &msg)
		fmt.Printf("recv data: %s\n", msg.String())
	})

	// resp ack
	ack, _ := codecIns.Encode([]byte("haha ack"))
	_, _ = c.Write(ack)
	return gnet.None
}
