package net

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/logging"
	"ppim/api/message_pb"
	"ppim/internal/comet/net/codec"
	"sync/atomic"
	"time"
)

type (
	EventEngine struct {
		gnet.BuiltinEventEngine
		context   *ServerContext
		processor *Processor
	}

	EventConnContext struct {
		Codec  codec.Codec
		Authed bool
	}
)

func newEventEngine(context *ServerContext) *EventEngine {
	return &EventEngine{
		context: context,
		processor: &Processor{
			context: context,
		},
	}
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

	atomic.AddInt32(&e.context.online, 1)
	return
}

func (e *EventEngine) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	if err != nil {
		logging.Infof("error occurred on connection=%s, %v\n", c.RemoteAddr().String(), err)
	}
	atomic.AddInt32(&e.context.online, -1)

	e.context.connManager.RemoveWithFD(c.Fd())
	return
}

func (e *EventEngine) OnTraffic(_c gnet.Conn) gnet.Action {
	connCtx := _c.Context().(*EventConnContext)

	buf, err := connCtx.Codec.Decode(_c)
	if err != nil {
		if errors.Is(err, codec.ErrIncompletePacket) {
			return gnet.None
		}
		logging.Errorf("failed to decode, %v", err)
		return gnet.Close
	}
	var msg message_pb.Message
	_ = proto.Unmarshal(buf, &msg)
	fmt.Printf("recv data: %s\n", msg.String())

	if !connCtx.Authed {
		if err = e.processor.Auth(_c, msg.GetConnectPacket()); err != nil {
			logging.Errorf("failed to auth the connection, %v", err)
			return gnet.Close
		}
	} else {
		if err = e.processor.Process(_c, &msg); err != nil {
			logging.Errorf("failed to process msg, %v", err)
		}
	}

	if _c.InboundBuffered() > 0 {
		if err = _c.Wake(nil); err != nil { // wake up the connection manually to avoid missing the leftover data
			logging.Errorf("failed to wake up the connection, %v", err)
			return gnet.Close
		}
	}
	return gnet.None
}

func (e *EventEngine) OnTick() (delay time.Duration, action gnet.Action) {
	return
}
