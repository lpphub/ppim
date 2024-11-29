package net

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/lpphub/golib/logger"
	"github.com/panjf2000/gnet/v2"
	"ppim/api/protocol"
	"ppim/internal/gate/net/codec"
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
		context:   context,
		processor: newProcessor(context),
	}
}

func (e *EventEngine) start() error {
	return gnet.Run(e, e.context.Addr, gnet.WithMulticore(true), gnet.WithReusePort(true),
		gnet.WithTicker(true))
}

func (e *EventEngine) OnBoot(_ gnet.Engine) gnet.Action {
	logger.Log().Info().Msgf("Listening and accepting TCP on %s", e.context.Addr)
	return gnet.None
}

func (e *EventEngine) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	//if remoteArr := c.RemoteAddr(); remoteArr != nil {
	//	IP := strings.Split(remoteArr.String(), ":")[0]
	//	// 可增加IP黑名单控制
	//	logger.Log().Info().Msgf("open new connection from %s", IP)
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
		logger.Log().Err(err).Msgf("error occurred on connection=%s", c.RemoteAddr().String())
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
		logger.Log().Err(err).Msg("failed to decode")
		return gnet.Close
	}
	var msg protocol.Message
	_ = proto.Unmarshal(buf, &msg)
	logger.Log().Debug().Msgf("recv data: %s", msg.String())

	if !connCtx.Authed {
		if err = e.processor.Auth(_c, msg.GetConnectPacket()); err != nil {
			logger.Log().Err(err).Msg("failed to auth the connection")
			return gnet.Close
		}
	} else {
		if err = e.processor.Process(_c, &msg); err != nil {
			logger.Log().Err(err).Msg("failed to process msg")
		}
	}

	if _c.InboundBuffered() > 0 {
		if err = _c.Wake(nil); err != nil { // wake up the connection manually to avoid missing the leftover data
			logger.Log().Err(err).Msg("failed to wake up the connection")
			return gnet.Close
		}
	}
	return gnet.None
}

func (e *EventEngine) OnTick() (delay time.Duration, action gnet.Action) {
	interval := time.Now().Add(-5 * time.Minute)
	cm := e.context.connManager
	for i, c := range cm.connMap {
		if interval.After(c.HeartbeatLastTime) { // 超过5分钟未收到心跳
			cm.RemoveWithFD(i)
		}
	}
	delay = 3 * time.Minute
	return
}
