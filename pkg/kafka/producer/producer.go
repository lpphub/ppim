package producer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	config *Config
	writer *kafka.Writer
	mtx    sync.RWMutex
	closed bool
}

func NewProducer(opts ...Option) (*Producer, error) {
	config := defaultConfig()

	for _, opt := range opts {
		opt(config)
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	writer := &kafka.Writer{
		Addr:         kafka.TCP(config.brokers...),
		Balancer:     &kafka.LeastBytes{},
		MaxAttempts:  config.maxRetries,
		RequiredAcks: config.requiredAcks,
		Async:        false,
		Logger:       config.logger,
	}
	if config.topic != "" {
		writer.Topic = config.topic
	}

	if config.clientID != "" {
		writer.Transport = &kafka.Transport{
			ClientID: config.clientID,
		}
	}

	return &Producer{
		config: config,
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
