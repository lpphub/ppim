package mq

import (
	"context"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/lpphub/golib/logger"
	"github.com/segmentio/kafka-go"
	"ppim/api/protocol"
	"ppim/internal/chatlib"
	"ppim/internal/gate/global"
	"ppim/internal/gate/net"
	"ppim/pkg/kafkago"
	"time"
)

type Subscriber struct {
	svc *net.ServerContext
}

func RegisterSubscriber(svc *net.ServerContext) error {
	sub := &Subscriber{
		svc: svc,
	}
	return sub.register()
}

func (s *Subscriber) register() error {
	kconf := global.Conf.Kafka

	config := kafkago.ConsumerConfig{
		Brokers:     kconf.Brokers,
		Topic:       kconf.Topic,
		GroupID:     kconf.GroupId,
		MaxWait:     time.Second,
		MaxAttempts: 3,
	}
	c, err := kafkago.NewConsumer(s.handleDelivery, config)
	if err != nil {
		return fmt.Errorf("register kafka consumer error: %v", err)
	}
	c.Start()
	return nil
}

func (s *Subscriber) handleDelivery(ctx context.Context, message kafka.Message) error {
	logger.Log().Info().Msgf("mq receive message: %s", string(message.Value))

	var msg chatlib.DeliveryMsg
	if err := jsoniter.Unmarshal(message.Value, &msg); err != nil {
		return err
	}

	switch msg.CMD {
	case chatlib.DeliveryChat:
		chat := msg.Chat
		bytes, _ := protocol.PacketReceive(&protocol.ReceivePacket{
			ConversationType: chat.ConversationType,
			ConversationID:   chat.ConversationID,
			FromUID:          chat.FromUID,
			Payload: &protocol.Payload{
				MsgNo:   chat.MsgNo,
				MsgId:   chat.MsgID,
				MsgSeq:  chat.MsgSeq,
				MsgType: int32(chat.MsgType),
				Content: chat.Content,
			},
		})

		for _, receiver := range msg.Receivers {
			clients := s.svc.ConnManager.GetWithUID(receiver)
			if len(clients) == 0 {
				continue
			}
			for _, client := range clients {
				if client.DID != chat.FromDID { // 避免自己收到消息
					_, err := client.Write(bytes)
					if err != nil {
						logger.Err(ctx, err, "write to client error")
					}

					// 放入重试队列
					s.svc.RetryManager.Add(&net.RetryMsg{
						MsgBytes: bytes,
						ConnFD:   client.Conn.Fd(),
						MsgID:    chat.MsgID,
						UID:      receiver,
					})
				}
			}
		}

	case chatlib.DeliveryEvent:
		// TODO: handle event
	case chatlib.DeliveryNotify:
		// TODO: handle notify
	default:
		logger.Warnf(ctx, "unknown cmd: %s", msg.CMD)
	}
	return nil
}
