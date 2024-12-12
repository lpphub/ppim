package consumer

import (
	"errors"
	"github.com/lpphub/golib/logger"
	"github.com/segmentio/kafka-go"
	"time"
)

type Option func(*Config)

type Config struct {
	brokers     []string
	topic       string
	groupID     string
	minBytes    int
	maxBytes    int
	maxAttempts int
	startOffset int64
	maxWait     time.Duration
	logger      *logger.Logger
}

func defaultConfig() *Config {
	return &Config{
		minBytes:    10e3, // 10KB
		maxBytes:    10e6, // 10MB
		maxWait:     3 * time.Second,
		maxAttempts: 3,
		startOffset: kafka.LastOffset,
		logger:      logger.Log(),
	}
}

func WithBrokers(brokers []string) Option {
	return func(c *Config) {
		c.brokers = brokers
	}
}

func WithTopic(topic string) Option {
	return func(c *Config) {
		c.topic = topic
	}
}

func WithGroupID(groupID string) Option {
	return func(c *Config) {
		c.groupID = groupID
	}
}

func WithMinBytes(minBytes int) Option {
	return func(c *Config) {
		c.minBytes = minBytes
	}
}

func WithMaxBytes(maxBytes int) Option {
	return func(c *Config) {
		c.maxBytes = maxBytes
	}
}

func WithMaxWait(maxWait time.Duration) Option {
	return func(c *Config) {
		c.maxWait = maxWait
	}
}

func WithMaxAttempts(maxAttempts int) Option {
	return func(c *Config) {
		c.maxAttempts = maxAttempts
	}
}

func WithStartOffset(offset int64) Option {
	return func(c *Config) {
		c.startOffset = offset
	}
}

func WithLogger(logger *logger.Logger) Option {
	return func(c *Config) {
		c.logger = logger
	}
}

func (c *Config) Validate() error {
	if len(c.brokers) == 0 || c.topic == "" || c.groupID == "" {
		return errors.New("invalid consumer config: brokers, topic and groupID are required")
	}
	return nil
}
