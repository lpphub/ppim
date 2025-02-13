package sub

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

func SubscribeDelivery(svc *net.ServerContext) error {
	s := &Subscriber{
		svc: svc,
	}
	return s.registerKafka()
}

func (s *Subscriber) registerKafka() error {
	fn := func(ctx context.Context, message kafka.Message) error {
		logger.Log().Info().Msgf("mq receive message: %s", string(message.Value))

		var msg chatlib.DeliveryMsg
		if err := jsoniter.Unmarshal(message.Value, &msg); err != nil {
			return err
		}
		s.delivery(ctx, &msg)
		return nil
	}

	kc, err := kafkago.NewConsumer(fn, kafkago.ConsumerConfig{
		Brokers:     global.Conf.Kafka.Brokers,
		Topic:       global.Conf.Kafka.Topic,
		GroupID:     global.Conf.Kafka.GroupId,
		MaxWait:     time.Second,
		MaxAttempts: 3,
	})
	if err != nil {
		return fmt.Errorf("register kafka consumer error: %v", err)
	}
	kc.Start()
	return nil
}

func (s *Subscriber) delivery(ctx context.Context, msg *chatlib.DeliveryMsg) {
	switch msg.CMD {
	case chatlib.DeliveryChat:
		s.deliveryChat(ctx, msg)
	case chatlib.DeliveryEvent:
		s.deliveryEvent(ctx, msg)
	case chatlib.DeliveryNotify:
		s.deliveryNotify(ctx, msg)
	default:
		logger.Warnf(ctx, "unknown cmd: %s", msg.CMD)
	}
}

func (s *Subscriber) deliveryChat(ctx context.Context, msg *chatlib.DeliveryMsg) {
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
}

func (s *Subscriber) deliveryEvent(ctx context.Context, msg *chatlib.DeliveryMsg) {
	// TODO: handle event
}

func (s *Subscriber) deliveryNotify(ctx context.Context, msg *chatlib.DeliveryMsg) {
	// TODO: handle notify
}
