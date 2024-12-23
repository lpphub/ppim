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

type ConsumerConfig struct {
	Brokers     []string
	Topic       string
	GroupID     string
	MinBytes    int
	MaxBytes    int
	MaxAttempts int
	StartOffset int64
	MaxWait     time.Duration
	Logger      *logger.Logger
}

func (c *ConsumerConfig) Validate() error {
	if len(c.Brokers) == 0 || c.Topic == "" || c.GroupID == "" {
		return errors.New("invalid consumer config: Brokers, topic and groupID are required")
	}
	return nil
}

type MessageHandler func(context.Context, kafka.Message) error

type Consumer struct {
	config   *ConsumerConfig
	reader   *kafka.Reader
	handler  MessageHandler
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
	stopOnce sync.Once
}

func NewConsumer(handler MessageHandler, config ConsumerConfig) (*Consumer, error) {
	if handler == nil {
		return nil, fmt.Errorf("message handler is required")
	}
	// default config
	if config.MinBytes == 0 {
		config.MinBytes = 10e3
	}
	if config.MaxBytes == 0 {
		config.MaxBytes = 10e6
	}
	if config.MaxWait <= 0 {
		config.MaxWait = 3 * time.Second
	}
	if config.StartOffset == 0 {
		config.StartOffset = kafka.LastOffset
	}
	if config.Logger == nil {
		config.Logger = logger.Log()
	}
	if err := config.Validate(); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	return &Consumer{
		config:  &config,
		handler: handler,
		ctx:     ctx,
		cancel:  cancel,
	}, nil
}

func (c *Consumer) Start() {
	c.config.Logger.Info().Msgf("start consumer: topic=%s, groupID=%s", c.config.Topic, c.config.GroupID)

	c.reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:     c.config.Brokers,
		Topic:       c.config.Topic,
		GroupID:     c.config.GroupID,
		MinBytes:    c.config.MinBytes,
		MaxBytes:    c.config.MaxBytes,
		MaxWait:     c.config.MaxWait,
		StartOffset: c.config.StartOffset,
		ErrorLogger: kafka.LoggerFunc(c.config.Logger.Printf),
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
					c.config.Logger.Printf("Error fetching message: %v", err)
				}
				continue
			}

			// 处理消息
			attempts := 0
			for attempts < c.config.MaxAttempts {
				err = c.handler(c.ctx, msg)
				if err == nil {
					break
				}
				attempts++
				c.config.Logger.Printf("Error processing message (attempt %d/%d): %v", attempts, c.config.MaxAttempts, err)
			}

			if err = c.reader.CommitMessages(c.ctx, msg); err != nil {
				c.config.Logger.Printf("Error committing message: %v", err)
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
