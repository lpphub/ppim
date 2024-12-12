package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"io"
	"os"
	"ppim/api/protocol"
	"ppim/internal/chatlib"
	"time"
)

func main() {
	c, _, err := websocket.DefaultDialer.Dial("ws://localhost:5051/", nil)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	defer c.Close()

	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				if errors.Is(err, io.EOF) {
					os.Exit(-1)
					return
				}
				fmt.Printf("%v\n", err)
				return
			}

			var msg protocol.Message
			_ = proto.Unmarshal(message, &msg)

			if msg.GetMsgType() == protocol.MsgType_CONNECT_ACK {
				fmt.Printf("连接结果：%d \n", msg.GetConnectAckPacket().GetCode())
			}

			if msg.GetMsgType() == protocol.MsgType_SEND_ACK {
				fmt.Printf("发送结果：code=%d msgId=%s \n", msg.GetSendAckPacket().GetCode(), msg.GetSendAckPacket().GetMsgId())
			}

			if msg.GetMsgType() == protocol.MsgType_RECEIVE {
				d := msg.GetReceivePacket()
				fmt.Printf("接收到的消息：data=%s fromID=%s convID=%s \n", d.GetPayload().GetContent(), d.GetFromID(), d.GetConversationID())
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
		fmt.Printf("write err: %v\n", err)
		return
	}

	// 模拟消息发送
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("输入消息: ")
		text, _ := reader.ReadString('\n')

		msg3, _ := protocol.PacketSend(&protocol.SendPacket{
			ConversationType: chatlib.ConvSingle,
			ToID:             "123",
			Payload: &protocol.Payload{
				MsgNo:    uuid.New().String(),
				MsgType:  1,
				Content:  text,
				SendTime: uint64(time.Now().UnixMilli()),
			},
		})
		err = c.WriteMessage(websocket.BinaryMessage, msg3)
		if err != nil {
			fmt.Printf("write err: %v\n", err)
			return
		}
	}
}
