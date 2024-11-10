package codec

type PacketType uint8

const (
	UNKNOWN     PacketType = iota //未知
	CONNECT                       //连接
	CONNECT_ACK                   //连接ack
	SEND                          //发送
	SEND_ACK                      //发送ack
	RECEIVE                       //接收
	RECEIVE_ACK                   //接收ack
	PING                          //心跳
	PONG                          //心跳ack
	DISCONNECT                    //断开连接
)

type Frame struct {
	MagicNumber uint16
	PacketType  PacketType
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
