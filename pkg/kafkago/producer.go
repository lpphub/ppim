package kafkago

import (
	"context"
	"errors"
	"fmt"
	"github.com/lpphub/golib/logger"
	"github.com/segmentio/kafka-go"
	"sync"
	"time"
)

type ProducerConfig struct {
	Brokers      []string
	Topic        string
	ClientID     string
	MaxAttempts  int
	RequiredAcks kafka.RequiredAcks
	BatchTimeout time.Duration
	Async        bool
}

func (c ProducerConfig) Validate() error {
	if len(c.Brokers) == 0 {
		return errors.New("invalid producer configuration")
	}
	return nil
}

type Producer struct {
	config *ProducerConfig
	writer *kafka.Writer
	mtx    sync.RWMutex
	closed bool
}

func NewProducer(config ProducerConfig) (*Producer, error) {
	if config.BatchTimeout <= 0 {
		config.BatchTimeout = 1 * time.Second
	}
	if config.RequiredAcks.String() == "unknown" {
		config.RequiredAcks = kafka.RequireNone
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	writer := &kafka.Writer{
		Addr:         kafka.TCP(config.Brokers...),
		Balancer:     &kafka.LeastBytes{},
		MaxAttempts:  config.MaxAttempts,
		RequiredAcks: config.RequiredAcks,
		BatchTimeout: config.BatchTimeout,
		Async:        config.Async,
		Logger:       logger.Log(),
	}
	if config.Topic != "" {
		writer.Topic = config.Topic
	}
	if config.ClientID != "" {
		writer.Transport = &kafka.Transport{
			ClientID: config.ClientID,
		}
	}

	return &Producer{
		config: &config,
		writer: writer,
	}, nil
}

func (p *Producer) SendMessage(ctx context.Context, messages ...kafka.Message) error {
	p.mtx.RLock()
	defer p.mtx.RUnlock()

	if p.closed {
		return errors.New("producer is closed")
	}

	err := p.writer.WriteMessages(ctx, messages...)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	return nil
}

func (p *Producer) Close() error {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	if p.closed {
		return nil
	}

	p.closed = true
	if err := p.writer.Close(); err != nil {
		return fmt.Errorf("failed to close producer: %w", err)
	}
	return nil
}
