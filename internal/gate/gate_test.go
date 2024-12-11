package gate

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/panjf2000/gnet/v2/pkg/logging"
	"net"
	"net/url"
	"os"
	"os/signal"
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

func TestWebsocket(t *testing.T) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: "localhost:5051", Path: "/"}
	logging.Infof("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		logging.Fatalf("dial: %s", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				logging.Infof("read err: %s", err)
				return
			}

			var msg protocol.Message
			_ = proto.Unmarshal(message, &msg)
			logging.Infof("recv: %v", msg.GetConnectAckPacket())
		}
	}()

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			msg, _ := proto.Marshal(&protocol.Message{
				MsgType: protocol.MsgType_CONNECT,
				Payload: &protocol.Message_ConnectPacket{
					ConnectPacket: &protocol.ConnectPacket{
						Uid:   "123",
						Did:   "abc",
						Token: "aaa",
					},
				},
			})
			err = c.WriteMessage(websocket.BinaryMessage, msg)
			if err != nil {
				logging.Infof("write err: %s", err)
				return
			}
		case <-interrupt:
			logging.Infof("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err = c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				logging.Infof("write close err: %s", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
