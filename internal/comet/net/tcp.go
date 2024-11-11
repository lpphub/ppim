package net

import (
	"context"
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

	Addr string

	connected   int32
	connManager *ClientManager
}

type ConnContext struct {
	Codec  codec.Codec
	Authed bool
}

func NewTCPServer(addr string) *TCPServer {
	return &TCPServer{
		Addr: fmt.Sprintf("tcp://%s", addr),
	}
}

func (s *TCPServer) Start() error {
	return gnet.Run(s, s.Addr, gnet.WithMulticore(true))
}

func (s *TCPServer) OnBoot(_ gnet.Engine) gnet.Action {
	logging.Infof("Listening and accepting TCP on %s\n", s.Addr)
	return gnet.None
}

func (s *TCPServer) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	//if remoteArr := c.RemoteAddr(); remoteArr != nil {
	//	IP := strings.Split(remoteArr.String(), ":")[0]
	//	// 可增加IP黑名单控制
	//	logging.Infof("open new connection from %s", IP)
	//}

	ctx := context.WithValue(context.Background(), "ctx", &ConnContext{
		Authed: false,
		Codec:  new(codec.ProtobufCodec),
	})
	c.SetContext(ctx)

	atomic.AddInt32(&s.connected, 1)
	return
}

func (s *TCPServer) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	if err != nil {
		logging.Infof("error occurred on connection=%s, %v\n", c.RemoteAddr().String(), err)
	}
	atomic.AddInt32(&s.connected, -1)

	s.connManager.RemoveWithFD(c.Fd())
	return
}

func (s *TCPServer) OnTraffic(c gnet.Conn) gnet.Action {
	connCtx := c.Context().(*ConnContext)
	codecIns := connCtx.Codec.(*codec.ProtobufCodec)

	buf, err := codecIns.Decode(c)
	if err != nil {
		if errors.Is(err, codec.ErrInvalidMagic) { //非法数据包关闭连接，不完整的包继续处理
			return gnet.Close
		}
		return gnet.None
	}

	if !connCtx.Authed {
		// todo 未授权，仅接受认证数据
	} else {
		// todo 已授权，处理业务
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
