package consumer

import (
	"errors"
	"github.com/lpphub/golib/logger"
	"github.com/segmentio/kafka-go"
)

// Option 定义配置选项函数类型
type Option func(*Config)

// Config 定义消费者配置
type Config struct {
	brokers     []string
	topic       string
	groupID     string
	minBytes    int
	maxBytes    int
	maxAttempts int
	startOffset int64
	logger      *logger.Logger
}

// defaultConfig 创建默认配置
func defaultConfig() *Config {
	return &Config{
		minBytes:    10e3, // 10KB
		maxBytes:    10e6, // 10MB
		maxAttempts: 3,
		startOffset: kafka.LastOffset,
		logger:      logger.Log(),
	}
}

// WithBrokers 设置 Kafka brokers
func WithBrokers(brokers []string) Option {
	return func(c *Config) {
		c.brokers = brokers
	}
}

// WithTopic 设置主题
func WithTopic(topic string) Option {
	return func(c *Config) {
		c.topic = topic
	}
}

// WithGroupID 设置消费者组ID
func WithGroupID(groupID string) Option {
	return func(c *Config) {
		c.groupID = groupID
	}
}

// WithMinBytes 设置最小字节数
func WithMinBytes(minBytes int) Option {
	return func(c *Config) {
		c.minBytes = minBytes
	}
}

// WithMaxBytes 设置最大字节数
func WithMaxBytes(maxBytes int) Option {
	return func(c *Config) {
		c.maxBytes = maxBytes
	}
}

// WithMaxAttempts 设置最大重试次数
func WithMaxAttempts(maxAttempts int) Option {
	return func(c *Config) {
		c.maxAttempts = maxAttempts
	}
}

// WithStartOffset 设置起始偏移量
func WithStartOffset(offset int64) Option {
	return func(c *Config) {
		c.startOffset = offset
	}
}

// WithLogger 设置日志记录器
func WithLogger(logger *logger.Logger) Option {
	return func(c *Config) {
		c.logger = logger
	}
}

// Validate 验证配置是否有效
func (c *Config) Validate() error {
	if len(c.brokers) == 0 || c.topic == "" || c.groupID == "" {
		return errors.New("invalid consumer config: brokers, topic and groupID are required")
	}
	return nil
}
