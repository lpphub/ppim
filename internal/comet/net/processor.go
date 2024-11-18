package net

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/panjf2000/gnet/v2"
	"ppim/api/message_pb"
)

type Processor struct {
	context *ServerContext
}

var (
	ErrAuthParamEmpty = errors.New("授权参数为空")
	ErrAuthFailure    = errors.New("授权校验失败")
)

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
	// todo 授权rpc接口
	authed := uid == "123" && token == "aaa"
	if !authed {
		return ErrAuthFailure
	}

	client := &Client{
		Conn: conn,
		UID:  uid,
		DID:  did,
	}
	_ = client.SetAuthResult(true)
	p.context.connManager.Add(client)

	//响应ack
	ack := &message_pb.Message{
		MsgType: message_pb.MsgType_CONNECT,
		Payload: &message_pb.Message_ConnectAckPacket{
			ConnectAckPacket: &message_pb.ConnectAckPacket{
				Ok: true,
			},
		},
	}
	data, _ := proto.Marshal(ack)
	if _, err := client.Write(data); err != nil {
		return err
	}
	return nil
}

func (p *Processor) Ping(c gnet.Conn, _ *message_pb.PingPacket) error {
	fmt.Printf("收到ping请求: %s\n", c.RemoteAddr().String())

	client := p.context.connManager.GetWithFD(c.Fd())
	if client != nil {
		pong := &message_pb.Message{
			MsgType: message_pb.MsgType_PONG,
			Payload: &message_pb.Message_PongPacket{
				PongPacket: &message_pb.PongPacket{},
			},
		}
		pongData, _ := proto.Marshal(pong)

		_, err := client.Write(pongData)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Processor) pushMsg() {
	//_ = goroutine.Default().Submit(func() {
	//	// todo 异步处理业务
	//
	//})
}
