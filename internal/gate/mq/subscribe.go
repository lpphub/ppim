package mq

import (
	"context"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/segmentio/kafka-go"
	"ppim/internal/chatlib"
	"ppim/internal/gate/global"
	"ppim/pkg/kafka/consumer"
)

func RegisterSubscriber() {
	var (
		brokers = global.Conf.Kafka.Brokers
	)
	err := subscriber(brokers, "test", "test", deliver)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
}

func subscriber(brokers []string, topic, groupID string, handler consumer.MessageHandler) error {
	c, err := consumer.NewConsumer(handler, consumer.WithBrokers(brokers), consumer.WithTopic(topic), consumer.WithGroupID(groupID))
	if err != nil {
		return err
	}
	c.Start()
	return nil
}

func deliver(ctx context.Context, message kafka.Message) error {
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
