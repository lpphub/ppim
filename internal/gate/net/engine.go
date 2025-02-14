package net

import (
	"context"
	"errors"
	"fmt"
	"github.com/lpphub/golib/logger"
	"github.com/panjf2000/gnet/v2"
	"ppim/api/protocol"
	"ppim/internal/gate/net/codec"
	"sync/atomic"
	"time"
)

type Network int8

const (
	_tcp Network = iota
	_ws
)

type (
	EngineOption struct {
		Network Network
		Addr    string
	}

	EventEngine struct {
		gnet.BuiltinEventEngine
		eng gnet.Engine

		opt EngineOption
		svc *ServerContext
	}

	EventConnContext struct {
		Authed  bool
		Network Network
		Codec   codec.Codec[codec.Frame]
		WsCodec *codec.WsCodec
	}
)

func newEventEngine(ctx *ServerContext, opt EngineOption) EventEngine {
	return EventEngine{
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
		Network: e.opt.Network,
		Authed:  false,
		Codec:   new(codec.FrameCodec),
		WsCodec: new(codec.WsCodec),
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
	var (
		connCtx  = _c.Context().(*EventConnContext)
		messages [][]byte
	)
	if connCtx.Network == _tcp {
		frames, err := connCtx.Codec.Decode(_c)
		if err != nil {
			if errors.Is(err, codec.ErrIncompletePacket) {
				return gnet.None
			}
			logger.Log().Err(err).Msg("failed to decode")
			return gnet.Close
		}
		if len(frames) == 0 {
			return gnet.None
		}
		for _, frame := range frames {
			messages = append(messages, frame.Body)
		}
	} else if connCtx.Network == _ws {
		ws := connCtx.WsCodec
		if ws.ReadBufferBytes(_c) == gnet.Close {
			return gnet.Close
		}
		wsMessages, err := ws.Decode(_c)
		if err != nil {
			return gnet.Close
		}
		if wsMessages == nil {
			return gnet.None
		}

		for _, wsm := range wsMessages {
			messages = append(messages, wsm.Payload)
		}
	}

	// process msg
	for _, msgBytes := range messages {
		msg, err := protocol.Unmarshal(msgBytes)
		if err != nil {
			logger.Log().Err(err).Msg("failed to unmarshal proto")
			continue
		}

		if !connCtx.Authed {
			if err = e.svc.processor.Auth(_c, msg.GetConnectPacket()); err != nil {
				logger.Log().Err(err).Msg("failed to auth the connection")
				return gnet.Close
			}
		} else {
			if err = e.svc.processor.Process(_c, msg); err != nil {
				logger.Log().Err(err).Msg("failed to process msg")
			}
		}
	}

	if _c.InboundBuffered() > 0 {
		if err := _c.Wake(nil); err != nil { // wake up the connection manually to avoid missing the leftover data
			logger.Log().Err(err).Msg("failed to wake up the connection")
			return gnet.Close
		}
	}
	return gnet.None
}

func (e *EventEngine) OnTick() (delay time.Duration, action gnet.Action) {
	delay = 3 * time.Minute

	if e.opt.Network == _ws { // tcp与ws 引用同一个connManager，只执行一个即可；后续可优化共用一个event_engine
		return
	}
	logger.Log().Info().Msgf("cleaning connections without heartbeat...")

	interval := time.Now().Add(-5 * time.Minute)
	cm := e.svc.ConnManager
	for i, c := range cm.connMap {
		if interval.After(c.HeartbeatLastTime) { // 超过5分钟未收到心跳
			cm.RemoveWithFD(i)
			logger.Log().Warn().Msgf("close the connection without heartbeat: uid=%s, did=%s", c.UID, c.DID)
		}
	}
	return
}
