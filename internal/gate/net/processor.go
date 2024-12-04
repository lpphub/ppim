package net

import (
	"context"
	"errors"
	"github.com/bwmarrin/snowflake"
	"github.com/golang/protobuf/proto"
	"github.com/lpphub/golib/logger"
	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/pool/goroutine"
	"github.com/spf13/cast"
	"ppim/api/protocol"
	"ppim/api/rpctypes"
	"ppim/internal/chat"
	"ppim/internal/gate/rpc"
	"time"
)

type Processor struct {
	context        *ServerContext
	workerPool     *goroutine.Pool
	msgIDGenerator *snowflake.Node
}

var (
	ErrAuthParamEmpty = errors.New("授权参数为空")
	ErrConnNotFound   = errors.New("连接不存在")
)

func newProcessor(context *ServerContext) *Processor {
	// todo 集群模式时需兼容
	node, _ := snowflake.NewNode(1)
	return &Processor{
		context:        context,
		msgIDGenerator: node,
		workerPool:     goroutine.Default(),
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
		ack, _ := proto.Marshal(protocol.PacketConnectAck(&protocol.ConnectAckPacket{
			Code: protocol.ConnAuthFail,
		}))

		connCtx := conn.Context().(*EventConnContext)
		ackBuf, err := connCtx.Codec.Encode(ack)
		if err != nil {
			logger.Err(ctx, err, "failed to decode")
			return err
		}
		if _, err = conn.Write(ackBuf); err != nil {
			logger.Err(ctx, err, "")
			return err
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
	p.context.connManager.Add(client)

	ack, _ := proto.Marshal(protocol.PacketConnectAck(&protocol.ConnectAckPacket{
		Code: protocol.OK,
	}))
	if _, err := client.Write(ack); err != nil {
		logger.Err(ctx, err, "")
		return err
	}
	return nil
}

func (p *Processor) Process(conn gnet.Conn, msg *protocol.Message) error {
	// 取得连接客户端
	client := p.context.connManager.GetWithFD(conn.Fd())
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
		logger.Log().Error().Msgf("process message error: %v", err)
	}
	return err
}

func (p *Processor) ping(_c *Client, _ *protocol.PingPacket) error {
	logger.Log().Debug().Msgf("UID=[%s] 收到ping请求", _c.UID)
	_c.HeartbeatLastTime = time.Now()

	pongData, _ := proto.Marshal(protocol.PacketPong(&protocol.PongPacket{}))
	_, err := _c.Write(pongData)
	return err
}

func (p *Processor) send(_c *Client, message *protocol.SendPacket) error {
	var (
		msgId             = p.msgIDGenerator.Generate().String()
		msgSeq            = cast.ToUint64(msgId) // todo 消息序列号（先暂用msgId简单替代, 后可建seq服务）
		conversationID, _ = chat.GenConversationID(_c.UID, message.ToID, message.ConversationType)
	)
	msg := &rpctypes.MessageReq{
		FromID:           _c.UID,
		ToID:             message.ToID,
		ConversationType: message.ConversationType,
		ConversationID:   conversationID,
		MsgID:            msgId,
		MsgSeq:           msgSeq,
		MsgNo:            message.Payload.MsgNo,
		MsgType:          message.Payload.MsgType,
		Content:          message.Payload.Content,
	}
	logger.Log().Debug().Msgf("UID=[%s]发送消息: %v", _c.UID, msg)

	// 异步处理业务
	err := p.workerPool.Submit(func() {
		ctx := logger.WithCtx(context.Background())

		var bytes []byte
		_, err := rpc.Caller().SendMsg(ctx, msg)
		if err != nil {
			logger.Err(ctx, err, "rpc - send msg error")

			bytes, _ = proto.Marshal(protocol.PacketSendAck(&protocol.SendAckPacket{
				Code: protocol.SendFail,
			}))
		} else {
			bytes, _ = proto.Marshal(protocol.PacketSendAck(&protocol.SendAckPacket{
				Code: protocol.OK,
			}))
		}

		if _, err = _c.Write(bytes); err != nil {
			logger.Err(ctx, err, "write ack error: send msg")
		}
	})
	return err
}

func (p *Processor) receiveAck(_c *Client, msg *protocol.ReceiveAckPacket) error {
	// todo 接收客户端收到消息的确认
	// 1. 更新消息状态
	// 2. 结束
	return nil
}
