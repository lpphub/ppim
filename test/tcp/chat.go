package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
	"io"
	"net"
	"os"
	"ppim/api/protocol"
	"ppim/internal/chatlib"
	"ppim/internal/gate/net/codec"
	"time"
)

func main() {
	c, err := net.Dial("tcp", ":5001")
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	codecInst := new(codec.ProtobufCodec)

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

					if msg.GetConnectAckPacket().GetCode() == 0 {
						go func() {
							ping, _ := protocol.PacketPing(&protocol.PingPacket{})
							pingBytes, _ := codecInst.Encode(ping)

							ticker := time.NewTicker(30 * time.Second)
							for range ticker.C {
								if _, err = c.Write(pingBytes); err != nil {
									fmt.Printf("%v\n", err)
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
				}

				if msg.GetMsgType() == protocol.MsgType_PONG {
					fmt.Printf("心跳消息：%v \n", time.Now())
				}
			}
		}
	}()

	// 1. 连接授权
	msg1, _ := protocol.PacketConnect(&protocol.ConnectPacket{
		Uid:   "123",
		Did:   "a123",
		Token: "aaa",
	})
	buf, _ := codecInst.Encode(msg1)

	if _, err = c.Write(buf); err != nil {
		fmt.Printf("%v\n", err)
	}

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
