package main

import (
	"encoding/binary"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"net"
	"ppim/api/protocol"
	"ppim/internal/chatlib"
	"ppim/internal/gate/net/codec"
	"testing"
	"time"
)

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

	msg1, _ := protocol.PacketConnect(&protocol.ConnectPacket{
		Uid:   "123",
		Did:   "a123",
		Token: "aaa",
	})
	buf, _ := codecInst.Encode(msg1)

	if _, err = c.Write(buf); err != nil {
		t.Log(err)
	}

	msg2, _ := protocol.PacketSend(&protocol.SendPacket{
		ConversationType: chatlib.ConvSingle,
		ToID:             "456",
		Payload: &protocol.Payload{
			MsgNo:    "u123",
			MsgType:  1,
			Content:  "hello world",
			SendTime: uint64(time.Now().UnixMilli()),
		},
	})
	buf1, _ := codecInst.Encode(msg2)
	if _, err = c.Write(buf1); err != nil {
		t.Log(err)
	}

	select {}
}

func TestClient_2(t *testing.T) {
	c, _, err := websocket.DefaultDialer.Dial("ws://localhost:5051/", nil)
	if err != nil {
		t.Log(err)
		return
	}
	defer c.Close()

	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				t.Log(err)
				return
			}

			var msg protocol.Message
			_ = proto.Unmarshal(message, &msg)

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
	}()

	// 1. 连接授权
	msg1, _ := protocol.PacketConnect(&protocol.ConnectPacket{
		Uid:   "456",
		Did:   "p456",
		Token: "bbb",
	})
	err = c.WriteMessage(websocket.BinaryMessage, msg1)
	if err != nil {
		t.Logf("write err: %v", err)
		return
	}

	// 2. 发送消息
	msg2, _ := protocol.PacketSend(&protocol.SendPacket{
		ConversationType: chatlib.ConvSingle,
		ToID:             "123",
		Payload: &protocol.Payload{
			MsgNo:    "u124",
			MsgType:  1,
			Content:  "你好",
			SendTime: uint64(time.Now().UnixMilli()),
		},
	})
	err = c.WriteMessage(websocket.BinaryMessage, msg2)
	if err != nil {
		t.Logf("write err: %v", err)
		return
	}

	select {}
}
