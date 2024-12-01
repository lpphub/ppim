package gate

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"net"
	"ppim/api/protocol"
	"ppim/internal/gate/net/codec"
	"testing"
	"time"
)

func TestServe(t *testing.T) {
	Serve()
}

func TestClient(t *testing.T) {
	c, err := net.Dial("tcp", ":5050")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	codecIns := new(codec.ProtobufCodec)

	msg := &protocol.Message{
		MsgType: protocol.MsgType_CONNECT,
		Payload: &protocol.Message_ConnectPacket{
			ConnectPacket: &protocol.ConnectPacket{
				Uid:   "123",
				Did:   "ios01",
				Token: "aaa",
			},
		},
	}
	pbData, _ := proto.Marshal(msg)

	buf, _ := codecIns.Encode(pbData)
	_, err = c.Write(buf)
	if err != nil {
		fmt.Println(err.Error())
	}

	ackBuf := make([]byte, 2048)
	_, _ = c.Read(ackBuf)
	b, _ := codecIns.Unpack(ackBuf)

	var ackMsg protocol.Message
	_ = proto.Unmarshal(b, &ackMsg)
	fmt.Printf("%v\n", &ackMsg)

	if ackMsg.GetConnectAckPacket().GetCode() == 0 {
		fmt.Println("connect success")

		for range 20 {
			fmt.Println("ping request...")
			ping := &protocol.Message{
				MsgType: protocol.MsgType_PING,
				Payload: &protocol.Message_PingPacket{
					PingPacket: &protocol.PingPacket{},
				},
			}

			pingData, _ := proto.Marshal(ping)

			pingBuf, _ := codecIns.Encode(pingData)
			_, _ = c.Write(pingBuf)

			time.Sleep(3 * time.Second)
		}

	}
}
