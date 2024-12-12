package mq

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	"github.com/lpphub/golib/logger"
	"github.com/segmentio/kafka-go"
	"ppim/api/protocol"
	"ppim/internal/chatlib"
	"ppim/internal/gate/global"
	"ppim/internal/gate/net"
	"ppim/pkg/kafka/consumer"
)

type Subscriber struct {
	svc *net.ServerContext
}

func RegisterSubscriber(svc *net.ServerContext) {
	sub := &Subscriber{svc: svc}
	sub.register()
}

func (s *Subscriber) register() {
	kconf := global.Conf.Kafka

	if c, err := consumer.NewConsumer(s.handleDeliver, consumer.WithBrokers(kconf.Brokers), consumer.WithTopic(kconf.Topic),
		consumer.WithGroupID(kconf.GroupId)); err != nil {
		logger.Log().Err(err).Msg("MQ subscriber register fail")
	} else {
		c.Start()
	}
}

func (s *Subscriber) handleDeliver(ctx context.Context, message kafka.Message) error {
	var msg chatlib.DeliverMsg
	if err := jsoniter.Unmarshal(message.Value, &msg); err != nil {
		return err
	}

	switch msg.CMD {
	case chatlib.DeliverChat:
		clients := s.svc.ConnManager.GetWithUID(msg.ToUID)
		if len(clients) == 0 {
			return nil
		}

		chat := msg.ChatMsg
		bytes, _ := protocol.PacketReceive(&protocol.ReceivePacket{
			ConversationType: chat.ConversationType,
			ConversationID:   chat.ConversationID,
			FromID:           chat.FromID,
			Payload: &protocol.Payload{
				MsgNo:   chat.MsgNo,
				MsgId:   chat.MsgID,
				MsgSeq:  chat.MsgSeq,
				MsgType: chat.MsgType,
				Content: chat.Content,
			},
		})

		for _, client := range clients {
			_, err := client.Write(bytes)
			if err != nil {
				logger.Err(ctx, err, "write to client error")
			}
		}
	case chatlib.DeliverEvent:
		// TODO: handle event
	case chatlib.DeliverNotify:
		// TODO: handle notify
	default:
		logger.Warnf(ctx, "unknown cmd: %s", msg.CMD)
	}
	return nil
}
