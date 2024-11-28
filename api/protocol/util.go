package protocol

const (
	OK           = iota
	ConnAuthFail // 连接认证失败
)

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
