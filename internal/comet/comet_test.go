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

	for i := range 5 {
		payload := &message_pb.ConnectPacket{
			UserId: fmt.Sprintf("uid-%d", i),
			Token:  "aaa",
		}

		msg := &message_pb.Message{
			MsgType: message_pb.MsgType_CONNECT,
			Payload: &message_pb.Message_ConnectPacket{
				ConnectPacket: payload,
			},
		}
		pbData, _ := proto.Marshal(msg)

		buf, _ := codecIns.Encode(pbData)
		_, err = c.Write(buf)
		if err != nil {
			fmt.Println(err.Error())
		}

		time.Sleep(1 * time.Second)
	}

	ack := make([]byte, 1024)
	_, _ = c.Read(ack)
	b, _ := codecIns.Unpack(ack)
	fmt.Println(string(b))
}
