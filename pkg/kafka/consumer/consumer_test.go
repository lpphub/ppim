package consumer

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	"github.com/lpphub/golib/logger"
	"github.com/segmentio/kafka-go"
	"ppim/internal/chatlib"
	"testing"
)

func TestConsumer_Start(t *testing.T) {
	logger.Setup()

	c, err := NewConsumer(func(_ context.Context, message kafka.Message) error {

		var dd chatlib.DeliverData
		if err := jsoniter.Unmarshal(message.Value, &dd); err != nil {
			return err
		}
		t.Logf("%v", dd)

		return nil
	}, WithBrokers([]string{"kafka:9092"}), WithTopic("test"), WithGroupID("test"))
	if err != nil {
		t.Logf("%v\n", err)
	}
	c.Start()

	select {}
}
