package protocol

const (
	OK           = iota
	ConnAuthFail // 连接认证失败
	SendFail     // 发送失败
)

func PacketConnect(conn *ConnectPacket) *Message {
	return &Message{
		MsgType: MsgType_CONNECT,
		Payload: &Message_ConnectPacket{
			ConnectPacket: conn,
		},
	}
}

func PacketConnectAck(connAck *ConnectAckPacket) *Message {
	return &Message{
		MsgType: MsgType_CONNECT_ACK,
		Payload: &Message_ConnectAckPacket{
			ConnectAckPacket: connAck,
		},
	}
}

func PacketPong(pong *PongPacket) *Message {
	return &Message{
		MsgType: MsgType_PONG,
		Payload: &Message_PongPacket{
			PongPacket: pong,
		},
	}
}

func PacketSendAck(sendAck *SendAckPacket) *Message {
	return &Message{
		MsgType: MsgType_SEND_ACK,
		Payload: &Message_SendAckPacket{
			SendAckPacket: sendAck,
		},
	}
}

func PacketReceive(receive *ReceivePacket) *Message {
	return &Message{
		MsgType: MsgType_RECEIVE,
		Payload: &Message_ReceivePacket{
			ReceivePacket: receive,
		},
	}
}
