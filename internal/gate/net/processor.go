package net

import (
	"context"
	"errors"
	"github.com/gobwas/ws/wsutil"
	"github.com/lpphub/golib/gowork"
	"github.com/lpphub/golib/logger"
	"github.com/panjf2000/gnet/v2"
	"github.com/rs/xid"
	"ppim/api/protocol"
	"ppim/internal/chatlib"
	"ppim/internal/gate/rpc"
	"time"
)

type Processor struct {
	svc        *ServerContext
	workerPool *gowork.Pool
}

var (
	ErrAuthParamEmpty = errors.New("授权参数为空")
	ErrConnNotFound   = errors.New("连接不存在")
)

func newProcessor(svc *ServerContext) *Processor {
	return &Processor{
		svc:        svc,
		workerPool: gowork.Default(),
	}
}

func (p *Processor) Auth(conn gnet.Conn, packet *protocol.ConnectPacket) error {
	var (
		uid   = packet.GetUid()
		did   = packet.GetDid()
		token = packet.GetToken()
	)
	if uid == "" || did == "" || token == "" {
		return ErrAuthParamEmpty
	}
	ctx := logger.WithCtx(context.Background())
	logger.Infof(ctx, "conn auth param: uid=%s, did=%s, token=%s", uid, did, token)

	authed, _ := rpc.Caller().Auth(ctx, uid, did, token)
	if !authed {
		ack, _ := protocol.PacketConnectAck(&protocol.ConnectAckPacket{
			Code: protocol.ConnAuthFail,
		})

		connCtx := conn.Context().(*EventConnContext)
		if connCtx.Network == _ws {
			if err := wsutil.WriteServerBinary(conn, ack); err != nil {
				logger.Err(ctx, err, "failed to write ws")
				return err
			}
		} else {
			ackBuf, err := connCtx.Codec.Encode(ack)
			if err != nil {
				logger.Err(ctx, err, "failed to decode")
				return err
			}
			if _, err = conn.Write(ackBuf); err != nil {
				logger.Err(ctx, err, "failed to write tcp")
				return err
			}
		}
		return nil
	}

	client := &Client{
		Conn:              conn,
		UID:               uid,
		DID:               did,
		HeartbeatLastTime: time.Now(),
	}
	_ = client.SetAuthResult(true)
	p.svc.ConnManager.Add(client)

	ack, _ := protocol.PacketConnectAck(&protocol.ConnectAckPacket{
		Code: protocol.OK,
	})
	if _, err := client.Write(ack); err != nil {
		logger.Err(ctx, err, "")
		return err
	}
	return nil
}

func (p *Processor) Process(conn gnet.Conn, msg *protocol.Message) error {
	// 取得连接客户端
	client := p.svc.ConnManager.GetWithFD(conn.Fd())
	if client == nil {
		return ErrConnNotFound
	}

	var err error
	switch msg.MsgType {
	case protocol.MsgType_PING:
		err = p.ping(client, msg.GetPingPacket())
	case protocol.MsgType_SEND:
		err = p.send(client, msg.GetSendPacket())
	case protocol.MsgType_RECEIVE_ACK:
		err = p.receiveAck(client, msg.GetReceiveAckPacket())
	default:
		err = errors.New("unknown message type")
	}
	if err != nil {
		logger.Log().Err(err).Msgf("process message error: %v", msg)
	}
	return err
}

func (p *Processor) ping(_c *Client, _ *protocol.PingPacket) error {
	logger.Log().Debug().Msgf("UID=[%s] 收到ping请求", _c.UID)
	// todo 心跳时，更新客户端route状态

	_c.HeartbeatLastTime = time.Now()

	bytes, _ := protocol.PacketPong(nil)
	_, err := _c.Write(bytes)
	return err
}

func (p *Processor) send(_c *Client, message *protocol.SendPacket) error {
	var (
		ctx               = logger.WithCtx(context.Background())
		conversationID, _ = chatlib.GenConversationID(message.ConversationType, _c.UID, message.ToID)
		msgId             = xid.New().String()
	)

	msg := &chatlib.MessageReq{
		FromUID:          _c.UID,
		FromDID:          _c.DID,
		ToID:             message.ToID,
		ConversationType: message.ConversationType,
		ConversationID:   conversationID,
		MsgID:            msgId,
		MsgNo:            message.Payload.MsgNo,
		MsgType:          int8(message.Payload.MsgType),
		Content:          message.Payload.Content,
		SendTime:         int64(message.Payload.SendTime),
		CreatedAt:        time.Now().UnixMilli(),
	}
	logger.Debugf(ctx, "UID=[%s]发送消息: %v", _c.UID, msg)

	// 异步处理业务
	err := p.workerPool.Submit(func() {
		var bytes []byte
		resp, serr := rpc.Caller().SendMsg(ctx, msg)
		if serr != nil {
			logger.Err(ctx, serr, "rpc: send msg error")

			bytes, _ = protocol.PacketSendAck(&protocol.SendAckPacket{
				Code:  protocol.SendFail,
				MsgNo: msg.MsgNo,
			})
		} else {
			bytes, _ = protocol.PacketSendAck(&protocol.SendAckPacket{
				Code:   protocol.OK,
				MsgNo:  msg.MsgNo,
				MsgId:  msg.MsgID,
				MsgSeq: resp.MsgSeq,
			})
		}

		if _, serr = _c.Write(bytes); serr != nil {
			logger.Err(ctx, serr, "write ack error: send msg")
		}
	})
	return err
}

func (p *Processor) receiveAck(_c *Client, msg *protocol.ReceiveAckPacket) error {
	logger.Infof(context.Background(), "receive ack msg: uid=%s, msgId=%s", _c.UID, msg.MsgId)

	err := p.workerPool.Submit(func() {
		// 从重试队列移除
		p.svc.RetryManager.Remove(_c.Conn.Fd(), msg.GetMsgId())
	})
	return err
}
