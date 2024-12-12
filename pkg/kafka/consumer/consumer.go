package consumer

import (
	"context"
	"errors"
	"fmt"
	"github.com/segmentio/kafka-go"
	"sync"
)

type MessageHandler func(context.Context, kafka.Message) error

type Consumer struct {
	config   *Config
	reader   *kafka.Reader
	handler  MessageHandler
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
	stopOnce sync.Once
}

func NewConsumer(handler MessageHandler, opts ...Option) (*Consumer, error) {
	if handler == nil {
		return nil, fmt.Errorf("message handler is required")
	}

	config := defaultConfig()

	for _, opt := range opts {
		opt(config)
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &Consumer{
		config:  config,
		handler: handler,
		ctx:     ctx,
		cancel:  cancel,
	}, nil
}

func (c *Consumer) Start() {
	c.config.logger.Info().Msgf("start consumer: topic=%s, groupID=%s", c.config.topic, c.config.groupID)

	c.reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:     c.config.brokers,
		Topic:       c.config.topic,
		GroupID:     c.config.groupID,
		MinBytes:    c.config.minBytes,
		MaxBytes:    c.config.maxBytes,
		StartOffset: c.config.startOffset,
		//ReadBackoffMin: 10 * time.Millisecond,
		//ReadBackoffMax: 200 * time.Millisecond,
		ErrorLogger: kafka.LoggerFunc(c.config.logger.Printf),
	})

	c.wg.Add(1)
	go c.consume()
}

func (c *Consumer) consume() {
	defer c.wg.Done()
	defer c.reader.Close()

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			msg, err := c.reader.FetchMessage(c.ctx)
			if err != nil {
				if !errors.Is(err, context.Canceled) {
					c.config.logger.Printf("Error fetching message: %v", err)
				}
				continue
			}

			// 处理消息
			attempts := 0
			for attempts < c.config.maxAttempts {
				err = c.handler(c.ctx, msg)
				if err == nil {
					break
				}
				attempts++
				c.config.logger.Printf("Error processing message (attempt %d/%d): %v", attempts, c.config.maxAttempts, err)
			}

			if err = c.reader.CommitMessages(c.ctx, msg); err != nil {
				c.config.logger.Printf("Error committing message: %v", err)
			}
		}
	}
}

func (c *Consumer) Stop() {
	c.stopOnce.Do(func() {
		c.cancel()
		c.wg.Wait()
	})
}
