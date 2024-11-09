package packet

type FrameType uint8

const (
	UNKNOWN     FrameType = iota //未知
	CONNECT                      //连接
	CONNECT_ACK                  //连接ack
	SEND                         //发送
	SEND_ACK                     //发送ack
	RECV                         //接收
	RECV_ACK                     //接收ack
	PING                         //心跳
	PONG                         //心跳ack
	DISCONNECT                   //断开连接
)

type Frame struct {
	MagicNumber uint16
	PacketType  FrameType
	Length      uint8
	Payload     []byte
}

type ConnectPacket struct {
	Frame
	UID string
}

type PingPacket struct {
	Frame
}

type PongPacket struct {
	Frame
}
