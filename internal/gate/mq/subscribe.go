package mq

import (
	"context"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/lpphub/golib/logger"
	"github.com/segmentio/kafka-go"
	"ppim/internal/chatlib"
	"ppim/internal/gate/global"
	"ppim/internal/gate/net"
	"ppim/pkg/kafka/consumer"
)

type Subscriber struct {
	svc *net.ServerContext
}

func LoadSubscriber(svc *net.ServerContext) {
	sub := &Subscriber{svc: svc}
	sub.register()
}

func (s *Subscriber) register() {
	brokers := global.Conf.Kafka.Brokers

	c, err := consumer.NewConsumer(s.deliver,
		consumer.WithBrokers(brokers),
		consumer.WithTopic("test"),
		consumer.WithGroupID("test"),
	)
	if err != nil {
		logger.Log().Err(err).Msg("MQ subscriber register fail")
	} else {
		c.Start()
	}
}

func (*Subscriber) deliver(ctx context.Context, message kafka.Message) error {
	var dd chatlib.DeliverData
	if err := jsoniter.Unmarshal(message.Value, &dd); err != nil {
		return err
	}

	for _, uid := range dd.ToUID {
		fmt.Println(uid)
		//clients := s.svc.ConnManager.GetWithUID(uid)
		//
		//for _, client := range clients {
		//	_, err := client.Write(msg.MsgData)
		//	if err != nil {
		//		logger.Err(ctx, err, "write to client error")
		//	}
		//}
	}
	return nil
}
