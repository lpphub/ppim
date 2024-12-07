package consumer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/segmentio/kafka-go"
)

// MessageHandler 定义消息处理函数类型
type MessageHandler func(context.Context, kafka.Message) error

// Consumer 定义消费者结构
type Consumer struct {
	config   *Config
	reader   *kafka.Reader
	handler  MessageHandler
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
	stopOnce sync.Once
}

// NewConsumer 创建新的消费者实例
func NewConsumer(handler MessageHandler, opts ...Option) (*Consumer, error) {
	if handler == nil {
		return nil, fmt.Errorf("message handler is required")
	}

	// 创建默认配置
	config := defaultConfig()

	// 应用所有选项
	for _, opt := range opts {
		opt(config)
	}

	// 验证配置
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

// Start 启动消费者
func (c *Consumer) Start() {
	c.reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:     c.config.brokers,
		Topic:       c.config.topic,
		GroupID:     c.config.groupID,
		MinBytes:    c.config.minBytes,
		MaxBytes:    c.config.maxBytes,
		StartOffset: c.config.startOffset,
		ErrorLogger: kafka.LoggerFunc(c.config.logger.Printf),
	})

	c.wg.Add(1)
	go c.consume()
}

// consume 消费消息的核心逻辑
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
				if !errors.Is(c.ctx.Err(), context.Canceled) {
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
				c.config.logger.Printf("Error processing message (attempt %d/%d): %v",
					attempts, c.config.maxAttempts, err)
			}

			if err := c.reader.CommitMessages(c.ctx, msg); err != nil {
				c.config.logger.Printf("Error committing message: %v", err)
			}
		}
	}
}

// Stop 优雅停止消费者
func (c *Consumer) Stop() {
	c.stopOnce.Do(func() {
		c.cancel()
		c.wg.Wait()
	})
}
