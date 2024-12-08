package producer

import (
	"errors"
	"github.com/lpphub/golib/logger"
	"github.com/segmentio/kafka-go"
)

type Option func(*Config)

type Config struct {
	brokers      []string
	topic        string
	clientID     string
	maxRetries   int
	requiredAcks kafka.RequiredAcks
	logger       *logger.Logger
}

func defaultConfig() *Config {
	c := &Config{
		maxRetries:   3,
		requiredAcks: kafka.RequireAll,
		logger:       logger.Log(),
	}
	return c
}

// WithBrokers sets the Kafka brokers
func WithBrokers(brokers []string) Option {
	return func(c *Config) {
		c.brokers = brokers
	}
}

// WithTopic sets the topic
func WithTopic(topic string) Option {
	return func(c *Config) {
		c.topic = topic
	}
}

// WithClientID sets the client ID
func WithClientID(clientID string) Option {
	return func(c *Config) {
		c.clientID = clientID
	}
}

// WithMaxRetries sets the maximum number of retries
func WithMaxRetries(maxRetries int) Option {
	return func(c *Config) {
		c.maxRetries = maxRetries
	}
}

// WithRequiredAcks sets the required acknowledgments
func WithRequiredAcks(requiredAcks kafka.RequiredAcks) Option {
	return func(c *Config) {
		c.requiredAcks = requiredAcks
	}
}

// WithLogger sets the logger
func WithLogger(logger *logger.Logger) Option {
	return func(c *Config) {
		c.logger = logger
	}
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if len(c.brokers) == 0 || c.topic == "" {
		return errors.New("invalid producer configuration")
	}
	return nil
}
