package producer

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	"github.com/lpphub/golib/logger"
	"github.com/segmentio/kafka-go"
	"ppim/internal/chatlib"
	"testing"
)

func TestProducer_SendMessage(t *testing.T) {
	logger.Setup()

	p, err := NewProducer(WithBrokers([]string{"kafka:9092"}))
	if err != nil {
		t.Logf("%v", err)
		return
	}

	b, _ := jsoniter.Marshal(chatlib.DeliverMsg{
		CMD:     "event",
		ToUID:   []string{"789"},
		Content: []byte("hello"),
	})

	msg := kafka.Message{
		Topic: "test",
		Key:   []byte("test"),
		Value: b,
	}

	err = p.SendMessage(context.Background(), msg)
	if err != nil {
		t.Logf("%v", err)
	}
}
