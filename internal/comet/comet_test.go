package comet

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"net"
	"ppim/api/message_pb"
	"ppim/internal/comet/net/codec"
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

	msg := &message_pb.Message{
		MsgType: message_pb.MsgType_CONNECT,
		Payload: &message_pb.Message_ConnectPacket{
			ConnectPacket: &message_pb.ConnectPacket{
				Uid:   "123",
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

	var ackMsg message_pb.Message
	ackBuf := make([]byte, 1024)
	_, _ = c.Read(ackBuf)

	b, _ := codecIns.Unpack(ackBuf)
	_ = proto.Unmarshal(b, &ackMsg)
	fmt.Printf("%v\n", &ackMsg)

	if ackMsg.GetConnectAckPacket().GetOk() {
		fmt.Println("connect success")

		for range 3 {
			ping := &message_pb.Message{
				MsgType: message_pb.MsgType_PING,
				Payload: &message_pb.Message_PingPacket{
					PingPacket: &message_pb.PingPacket{},
				},
			}

			pingData, _ := proto.Marshal(ping)

			pingBuf, _ := codecIns.Encode(pingData)
			_, _ = c.Write(pingBuf)

			time.Sleep(3 * time.Second)
		}

	}
}
