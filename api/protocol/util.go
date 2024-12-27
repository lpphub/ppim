package protocol

import (
	"google.golang.org/protobuf/proto"
)

const (
	OK           = iota
	ConnAuthFail // 连接认证失败
	SendFail     // 发送失败
)

func PacketConnect(conn *ConnectPacket) ([]byte, error) {
	return proto.Marshal(&Message{
		MsgType: MsgType_CONNECT,
		Payload: &Message_ConnectPacket{
			ConnectPacket: conn,
		},
	})
}

func PacketConnectAck(connAck *ConnectAckPacket) ([]byte, error) {
	return proto.Marshal(&Message{
		MsgType: MsgType_CONNECT_ACK,
		Payload: &Message_ConnectAckPacket{
			ConnectAckPacket: connAck,
		},
	})
}

func PacketPing(ping *PingPacket) ([]byte, error) {
	return proto.Marshal(&Message{
		MsgType: MsgType_PING,
		Payload: &Message_PingPacket{
			PingPacket: ping,
		},
	})
}

func PacketPong(pong *PongPacket) ([]byte, error) {
	return proto.Marshal(&Message{
		MsgType: MsgType_PONG,
		Payload: &Message_PongPacket{
			PongPacket: pong,
		},
	})
}

func PacketSend(send *SendPacket) ([]byte, error) {
	return proto.Marshal(&Message{
		MsgType: MsgType_SEND,
		Payload: &Message_SendPacket{
			SendPacket: send,
		},
	})
}

func PacketSendAck(sendAck *SendAckPacket) ([]byte, error) {
	return proto.Marshal(&Message{
		MsgType: MsgType_SEND_ACK,
		Payload: &Message_SendAckPacket{
			SendAckPacket: sendAck,
		},
	})
}

func PacketReceive(receive *ReceivePacket) ([]byte, error) {
	return proto.Marshal(&Message{
		MsgType: MsgType_RECEIVE,
		Payload: &Message_ReceivePacket{
			ReceivePacket: receive,
		},
	})
}

func PacketReceiveAck(receiveAck *ReceiveAckPacket) ([]byte, error) {
	return proto.Marshal(&Message{
		MsgType: MsgType_RECEIVE_ACK,
		Payload: &Message_ReceiveAckPacket{
			ReceiveAckPacket: receiveAck,
		},
	})
}

func Unmarshal(buf []byte) (*Message, error) {
	var msg Message
	if err := proto.Unmarshal(buf, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}
