package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	"io"
	"os"
	"ppim/api/protocol"
	"ppim/internal/chatlib"
	"time"
)

func main() {
	c, _, _, err := ws.Dial(context.Background(), "ws://localhost:5002")
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	defer c.Close()

	go func() {
		for {
			message, _, err := wsutil.ReadServerData(c)
			if err != nil {
				if errors.Is(err, io.EOF) {
					os.Exit(-1)
					return
				}
				fmt.Printf("read err: %v\n", err)
				return
			}

			var msg protocol.Message
			_ = proto.Unmarshal(message, &msg)

			if msg.GetMsgType() == protocol.MsgType_CONNECT_ACK {
				fmt.Printf("连接结果：%d \n", msg.GetConnectAckPacket().GetCode())

				if msg.GetConnectAckPacket().GetCode() == 0 {
					go func() {
						ping, _ := protocol.PacketPing(&protocol.PingPacket{})

						ticker := time.NewTicker(30 * time.Second)
						for range ticker.C {
							if perr := wsutil.WriteClientBinary(c, ping); perr != nil {
								fmt.Printf("%v\n", perr)
							}
						}
					}()
				}
			}

			if msg.GetMsgType() == protocol.MsgType_SEND_ACK {
				d := msg.GetSendAckPacket()
				fmt.Printf("发送结果：code=%d msgId=%s msgSeq=%d\n", d.GetCode(), d.GetMsgId(), d.GetMsgSeq())
			}

			if msg.GetMsgType() == protocol.MsgType_RECEIVE {
				d := msg.GetReceivePacket()
				fmt.Printf("接收到的消息：data=%s fromID=%s convID=%s \n", d.GetPayload().GetContent(), d.GetFromUID(), d.GetConversationID())

				// ack
				bytes, _ := protocol.PacketReceiveAck(&protocol.ReceiveAckPacket{
					MsgId:  d.GetPayload().GetMsgId(),
					MsgSeq: d.GetPayload().GetMsgSeq(),
				})
				_ = wsutil.WriteClientBinary(c, bytes)
			}

			if msg.GetMsgType() == protocol.MsgType_PONG {
				fmt.Printf("心跳消息：%v \n", time.Now())
			}
		}
	}()

	// 1. 连接授权
	msg1, _ := protocol.PacketConnect(&protocol.ConnectPacket{
		Uid:   "789",
		Did:   "p789",
		Token: "ccc",
	})
	err = wsutil.WriteClientBinary(c, msg1)
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
			ConversationType: chatlib.ConvGroup,
			ToID:             "100",
			Payload: &protocol.Payload{
				MsgNo:    uuid.New().String(),
				MsgType:  1,
				Content:  text,
				SendTime: uint64(time.Now().UnixMilli()),
			},
		})
		err = wsutil.WriteClientBinary(c, msg3)
		if err != nil {
			fmt.Printf("write err: %v\n", err)
			return
		}
	}
}
