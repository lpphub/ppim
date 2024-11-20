package message_pb

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
