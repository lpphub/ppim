package net

import (
	"context"
	"errors"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/golang/protobuf/proto"
	"github.com/lpphub/golib/logger"
	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/pool/goroutine"
	"ppim/api/message_pb"
	"ppim/internal/comet/rpc"
	"time"
)

type Processor struct {
	context        *ServerContext
	msgIDGenerator *snowflake.Node
}

var (
	ErrAuthParamEmpty = errors.New("授权参数为空")
	ErrAuthFailure    = errors.New("授权校验失败")
	ErrConnNotFound   = errors.New("连接不存在")
)

func newProcessor(context *ServerContext) *Processor {
	// todo 集群模式时需兼容
	node, _ := snowflake.NewNode(1)
	return &Processor{
		context:        context,
		msgIDGenerator: node,
	}
}

func (p *Processor) Auth(conn gnet.Conn, packet *message_pb.ConnectPacket) error {
	var (
		uid   = packet.GetUid()
		did   = packet.GetDid()
		token = packet.GetToken()
	)
	if uid == "" || did == "" || token == "" {
		return ErrAuthParamEmpty
	}
	fmt.Printf("auth param: uid=%s, did=%s, token=%s\n", uid, did, token)

	authed, _ := rpc.Caller().AuthSrv.Auth(logger.WithCtx(context.TODO()), uid, did, token)
	if !authed {
		return ErrAuthFailure
	}

	client := &Client{
		Conn:              conn,
		UID:               uid,
		DID:               did,
		HeartbeatLastTime: time.Now(),
	}
	_ = client.SetAuthResult(true)
	p.context.connManager.Add(client)

	// connect ack
	ack := message_pb.PacketConnectAck(&message_pb.ConnectAckPacket{
		Ok: true,
	})
	data, _ := proto.Marshal(ack)
	if _, err := client.Write(data); err != nil {
		return err
	}
	return nil
}

func (p *Processor) Process(conn gnet.Conn, msg *message_pb.Message) error {
	// 取得连接客户端
	client := p.context.connManager.GetWithFD(conn.Fd())
	if client == nil {
		return ErrConnNotFound
	}

	// 异步处理业务
	err := goroutine.Default().Submit(func() {
		// 消息分发
		var err error
		switch msg.MsgType {
		case message_pb.MsgType_PING:
			err = p.ping(client, msg.GetPingPacket())
		case message_pb.MsgType_SEND:
			err = p.send(client, msg.GetSendPacket())
		case message_pb.MsgType_RECEIVE_ACK:
			err = p.receiveAck(client, msg.GetReceiveAckPacket())
		default:
			err = errors.New("unknown message type")
		}
		if err != nil {
			fmt.Printf("process message error: %v\n", err)
		}
	})
	return err
}

func (p *Processor) ping(_c *Client, _ *message_pb.PingPacket) error {
	fmt.Printf("UID=[%s] 收到ping请求\n", _c.UID)
	_c.HeartbeatLastTime = time.Now()

	pongData, _ := proto.Marshal(message_pb.PacketPong(&message_pb.PongPacket{}))

	_, err := _c.Write(pongData)
	if err != nil {
		return err
	}
	return nil
}

func (p *Processor) send(_c *Client, msg *message_pb.SendPacket) error {
	// todo 接收客户端投递的消息
	// 1. 消息存储
	// 2. 消息在线投递
	// 3. 响应ack
	return nil
}

func (p *Processor) receiveAck(_c *Client, msg *message_pb.ReceiveAckPacket) error {
	// todo 接收客户端收到消息的确认
	// 1. 更新消息状态
	// 2. 结束
	return nil
}
