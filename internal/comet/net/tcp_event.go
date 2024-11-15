package net

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/logging"
	"github.com/panjf2000/gnet/v2/pkg/pool/goroutine"
	"ppim/api/message_pb"
	"ppim/internal/comet/net/codec"
	"sync/atomic"
)

type EventEngine struct {
	gnet.BuiltinEventEngine
	context *ServerContext
}

type EventConnContext struct {
	Codec  codec.Codec
	Authed bool
}

func newEventEngine(context *ServerContext) *EventEngine {
	return &EventEngine{context: context}
}

func (e *EventEngine) start() error {
	return gnet.Run(e, e.context.Addr, gnet.WithMulticore(true))
}

func (e *EventEngine) OnBoot(_ gnet.Engine) gnet.Action {
	logging.Infof("Listening and accepting TCP on %s\n", e.context.Addr)
	return gnet.None
}

func (e *EventEngine) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	//if remoteArr := c.RemoteAddr(); remoteArr != nil {
	//	IP := strings.Split(remoteArr.String(), ":")[0]
	//	// 可增加IP黑名单控制
	//	logging.Infof("open new connection from %s", IP)
	//}

	ctx := &EventConnContext{
		Authed: false,
		Codec:  new(codec.ProtobufCodec),
	}
	c.SetContext(ctx)

	atomic.AddInt32(&e.context.Connected, 1)
	return
}

func (e *EventEngine) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	if err != nil {
		logging.Infof("error occurred on connection=%s, %v\n", c.RemoteAddr().String(), err)
	}
	atomic.AddInt32(&e.context.Connected, -1)

	e.context.connManager.RemoveWithFD(c.Fd())
	return
}

func (e *EventEngine) OnTraffic(c gnet.Conn) gnet.Action {
	connCtx := c.Context().(*EventConnContext)

	buf, err := connCtx.Codec.Decode(c)
	if err != nil {
		if errors.Is(err, codec.ErrIncompletePacket) {
			return gnet.None
		}
		logging.Errorf("failed to decode, %v", err)
		return gnet.Close
	}

	if !connCtx.Authed {
		// todo 未授权时，进行授权校验，未通过则关闭连接
		connCtx.Authed = true
	} else {
		// todo 接收消息，处理业务逻辑
		fmt.Println("第二次")
	}

	_ = goroutine.Default().Submit(func() {
		// todo 异步处理业务

		var msg message_pb.Message
		_ = proto.Unmarshal(buf, &msg)
		fmt.Printf("recv data: %s\n", msg.String())
	})
	// resp ack
	ack, _ := connCtx.Codec.Encode([]byte("haha ack"))
	_, _ = c.Write(ack)

	if c.InboundBuffered() > 0 {
		if err = c.Wake(nil); err != nil { // wake up the connection manually to avoid missing the leftover data
			logging.Errorf("failed to wake up the connection, %v", err)
			return gnet.Close
		}
	}
	return gnet.None
}
