package net

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/lpphub/golib/logger"
	"github.com/panjf2000/gnet/v2"
	"ppim/api/protocol"
	"ppim/internal/gate/net/codec"
	"sync/atomic"
	"time"
)

type ConnType int8

const (
	_tcp ConnType = iota
	_ws
)

type (
	EngineOption struct {
		Protocol ConnType
		Addr     string
	}

	EventEngine struct {
		gnet.BuiltinEventEngine
		eng gnet.Engine

		opt EngineOption
		svc *ServerContext
	}

	EventConnContext struct {
		Codec    codec.Codec
		Authed   bool
		ConnType ConnType
		WsCodec  *codec.WsCodec
	}
)

func NewEventEngine(ctx *ServerContext, opt EngineOption) *EventEngine {
	return &EventEngine{
		svc: ctx,
		opt: opt,
	}
}

func (e *EventEngine) start() error {
	return gnet.Run(e, fmt.Sprintf("tcp://%s", e.opt.Addr),
		gnet.WithMulticore(true),
		gnet.WithReusePort(true),
		gnet.WithTicker(true),
	)
}

func (e *EventEngine) stop() error {
	return e.eng.Stop(context.Background())
}

func (e *EventEngine) OnBoot(eng gnet.Engine) gnet.Action {
	logger.Log().Info().Msgf("Listening and accepting Connection on %s", e.opt.Addr)
	e.eng = eng
	return gnet.None
}

func (e *EventEngine) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	//if remoteArr := c.RemoteAddr(); remoteArr != nil {
	//	IP := strings.Split(remoteArr.String(), ":")[0]
	//	// 可增加IP黑名单控制
	//	logger.Log().Info().Msgf("open new connection from %s", IP)
	//}

	ctx := &EventConnContext{
		ConnType: e.opt.Protocol,
		Authed:   false,
		Codec:    new(codec.ProtobufCodec), // todo 优化codec设计
		WsCodec:  new(codec.WsCodec),
	}
	c.SetContext(ctx)

	atomic.AddInt32(&e.svc.online, 1)
	return
}

func (e *EventEngine) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	if err != nil {
		logger.Log().Err(err).Msgf("error occurred on connection=%s", c.RemoteAddr().String())
	}
	atomic.AddInt32(&e.svc.online, -1)

	e.svc.ConnManager.RemoveWithFD(c.Fd())
	return
}

func (e *EventEngine) OnTraffic(_c gnet.Conn) gnet.Action {
	connCtx := _c.Context().(*EventConnContext)
	if connCtx.ConnType == _tcp {
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

		if act := e.process(_c, &msg, connCtx.Authed); act == gnet.Close {
			return act
		}
	} else if connCtx.ConnType == _ws {
		ws := connCtx.WsCodec
		if ws.ReadBufferBytes(_c) == gnet.Close {
			return gnet.Close
		}

		messages, err := ws.Decode(_c)
		if err != nil {
			return gnet.Close
		}
		if messages == nil {
			return gnet.None
		}

		for i := range messages {
			var msg protocol.Message
			_ = proto.Unmarshal(messages[i].Payload, &msg)

			if act := e.process(_c, &msg, connCtx.Authed); act == gnet.Close {
				return act
			}
		}
	}

	//if !connCtx.Authed {
	//	if err := e.svc.processor.Auth(_c, msg.GetConnectPacket()); err != nil {
	//		logger.Log().Err(err).Msg("failed to auth the connection")
	//		return gnet.Close
	//	}
	//} else {
	//	if err := e.svc.processor.Process(_c, &msg); err != nil {
	//		logger.Log().Err(err).Msg("failed to process msg")
	//	}
	//}

	if _c.InboundBuffered() > 0 {
		if err := _c.Wake(nil); err != nil { // wake up the connection manually to avoid missing the leftover data
			logger.Log().Err(err).Msg("failed to wake up the connection")
			return gnet.Close
		}
	}
	return gnet.None
}

func (e *EventEngine) process(_c gnet.Conn, msg *protocol.Message, isAuthed bool) gnet.Action {
	if !isAuthed {
		if err := e.svc.processor.Auth(_c, msg.GetConnectPacket()); err != nil {
			logger.Log().Err(err).Msg("failed to auth the connection")
			return gnet.Close
		}
	} else {
		if err := e.svc.processor.Process(_c, msg); err != nil {
			logger.Log().Err(err).Msg("failed to process msg")
		}
	}
	return gnet.None
}

func (e *EventEngine) OnTick() (delay time.Duration, action gnet.Action) {
	delay = 3 * time.Minute

	if e.opt.Protocol == _ws { // tcp与ws 引用同一个connManager，只执行一个即可；后续可优化共用一个event_engine
		return
	}
	logger.Log().Info().Msgf("cleaning expired connections...")

	interval := time.Now().Add(-5 * time.Minute)
	cm := e.svc.ConnManager
	for i, c := range cm.connMap {
		if interval.After(c.HeartbeatLastTime) { // 超过5分钟未收到心跳
			cm.RemoveWithFD(i)
		}
	}
	return
}