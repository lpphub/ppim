package sub

import (
	"context"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/segmentio/kafka-go"
	"ppim/internal/chatlib"
	"ppim/internal/gate/net"
	"ppim/pkg/kafka/consumer"
)

type Subscriber struct {
	svc *net.ServerContext
}

func subscriber(topic string) {
	c, err := consumer.NewConsumer(deliver,
		consumer.WithBrokers([]string{""}),
		consumer.WithTopic(topic),
		consumer.WithGroupID("test"),
	)
	if err != nil {
		panic(err)
	}
	c.Start()
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
