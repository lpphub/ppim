package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"io"
	"net"
	"os"
	"ppim/api/protocol"
	"ppim/internal/chatlib"
	"ppim/internal/gate/net/codec"
	"time"
)

func main() {
	c, err := net.Dial("tcp", ":5050")
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	go func() {
		for {
			headerLen := make([]byte, 4)
			if _, err = c.Read(headerLen); err != nil {
				if errors.Is(err, io.EOF) {
					os.Exit(-1)
					return
				}
				fmt.Printf("%v\n", err)
				return
			}

			payloadLen := binary.BigEndian.Uint32(headerLen)

			buf := make([]byte, payloadLen)
			if _, err := c.Read(buf); err != nil {
				fmt.Printf("%v\n", err)
				return
			} else {
				var msg protocol.Message
				_ = proto.Unmarshal(buf, &msg)

				if msg.GetMsgType() == protocol.MsgType_CONNECT_ACK {
					fmt.Printf("连接建立结果：%d \n", msg.GetConnectAckPacket().GetCode())
				}

				if msg.GetMsgType() == protocol.MsgType_SEND_ACK {
					fmt.Printf("发送结果：code=%d msgId=%s \n", msg.GetSendAckPacket().GetCode(), msg.GetSendAckPacket().GetMsgId())
				}

				if msg.GetMsgType() == protocol.MsgType_RECEIVE {
					d := msg.GetReceivePacket()
					fmt.Printf("接收到的消息：data=%s fromID=%s convID=%s \n", d.GetPayload().GetContent(), d.GetFromID(), d.GetConversationID())
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
		fmt.Printf("%v\n", err)
	}

	//msg2, _ := protocol.PacketSend(&protocol.SendPacket{
	//	ConversationType: chatlib.ConvSingle,
	//	ToID:             "456",
	//	Payload: &protocol.Payload{
	//		MsgNo:   "u123",
	//		MsgType: 1,
	//		Content: "hello world",
	//	},
	//})
	//buf1, _ := codecInst.Encode(msg2)
	//if _, err = c.Write(buf1); err != nil {
	//	fmt.Printf("%v\n", err)
	//}

	// 模拟消息发送
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("输入消息: ")
		text, _ := reader.ReadString('\n')

		msg3, _ := protocol.PacketSend(&protocol.SendPacket{
			ConversationType: chatlib.ConvSingle,
			ToID:             "456",
			Payload: &protocol.Payload{
				MsgNo:    uuid.New().String(),
				MsgType:  1,
				Content:  text,
				SendTime: uint64(time.Now().UnixMilli()),
			},
		})
		buf2, _ := codecInst.Encode(msg3)
		if _, err = c.Write(buf2); err != nil {
			fmt.Printf("%v\n", err)
		}
	}

}