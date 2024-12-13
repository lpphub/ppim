package producer

import (
	"context"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/lpphub/golib/logger"
	"github.com/segmentio/kafka-go"
	"ppim/internal/chatlib"
	"testing"
	"time"
)

func TestProducer_SendMessage(t *testing.T) {
	logger.Setup()

	p, err := NewProducer(WithBrokers([]string{"localhost:9094"}), WithBatchTimeout(100*time.Millisecond), WithAsync(false))
	if err != nil {
		t.Logf("%v", err)
		return
	}

	for i := 0; i < 3; i++ {
		b, _ := jsoniter.Marshal(chatlib.DeliverMsg{
			CMD:   "event_test",
			ToUID: fmt.Sprintf("u_%d", i),
		})
		msg := kafka.Message{
			Topic: "test_speed",
			Key:   []byte("test1"),
			Value: b,
		}

		t1 := time.Now()
		err = p.SendMessage(context.Background(), msg)
		if err != nil {
			t.Logf("%v", err)
		}
		t.Logf("cost %d", time.Since(t1).Milliseconds())
	}
}
