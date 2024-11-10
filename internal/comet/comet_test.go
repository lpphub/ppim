package comet

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"net"
	"ppim/api/message_pb"
	"ppim/internal/comet/net/codec"
	"testing"
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

	payload := &message_pb.ConnectPacket{
		UserId: "123",
		Token:  "bbb",
	}

	msg := &message_pb.Message{
		MsgType: message_pb.MsgType_CONNECT,
		Payload: &message_pb.Message_ConnectPacket{
			ConnectPacket: payload,
		},
	}
	pbData, _ := proto.Marshal(msg)

	codecIns := new(codec.ProtobufCodec)
	buf, _ := codecIns.Encode(pbData)
	//buf, _ := codecIns.Encode([]byte("hello tcp 3222"))
	_, err = c.Write(buf)
	if err != nil {
		fmt.Println(err.Error())
	}

	ack := make([]byte, 1024)
	_, _ = c.Read(ack)
	b, _ := codecIns.Unpack(ack)
	fmt.Println(string(b))
}
