package test

import (
	"encoding/binary"
	"github.com/golang/protobuf/proto"
	"github.com/lpphub/golib/env"
	"net"
	"ppim/api/protocol"
	"ppim/internal/gate"
	"ppim/internal/gate/net/codec"
	"ppim/internal/logic"
	"testing"
)

func TestGateServer(t *testing.T) {
	env.SetRootPath("../")
	gate.Serve()
}

func TestLogicServer(t *testing.T) {
	env.SetRootPath("../")
	logic.Serve()
}

func TestClient_1(t *testing.T) {
	c, err := net.Dial("tcp", ":5050")
	if err != nil {
		t.Log(err)
		return
	}

	go func() {
		for {
			headerLen := make([]byte, 4)
			if _, err = c.Read(headerLen); err != nil {
				t.Log(err)
				return
			}

			payloadLen := binary.BigEndian.Uint32(headerLen)

			buf := make([]byte, payloadLen)
			if _, err := c.Read(buf); err != nil {
				t.Log(err)
				return
			} else {
				var msg protocol.Message
				_ = proto.Unmarshal(buf, &msg)

				if msg.GetMsgType() == protocol.MsgType_CONNECT_ACK {
					t.Logf("recv: type=%v data=%v", msg.GetMsgType(), msg.GetConnectAckPacket())
				}

				if msg.GetMsgType() == protocol.MsgType_SEND_ACK {
					t.Logf("recv: type=%v data=%v", msg.GetMsgType(), msg.GetSendAckPacket())
				}

				if msg.GetMsgType() == protocol.MsgType_RECEIVE {
					t.Logf("recv: type=%v data=%v", msg.GetMsgType(), msg.GetReceivePacket())
				}
			}
		}
	}()

	// 1. 连接授权
	codecInst := new(codec.ProtobufCodec)

	msg := protocol.PacketConnect(&protocol.ConnectPacket{
		Uid:   "123",
		Did:   "ios01",
		Token: "aaa",
	})
	data, _ := proto.Marshal(msg)
	buf, _ := codecInst.Encode(data)

	if _, err = c.Write(buf); err != nil {
		t.Log(err)
	}

	select {}
}
