package protocol

import "github.com/golang/protobuf/proto"

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
