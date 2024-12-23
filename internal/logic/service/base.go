package service

import (
	"github.com/lpphub/golib/logger"
	"github.com/segmentio/kafka-go"
	"ppim/internal/logic/global"
	"ppim/pkg/kafkago"
	"time"
)

type ServiceContext struct {
	Msg   *MessageSrv
	Route *RouteSrv
	Conv  *ConversationSrv
}

var svc *ServiceContext

func InitService() {
	// kafka flush msg every 50ms
	conf := kafkago.ProducerConfig{
		Brokers:      global.Conf.Kafka.Brokers,
		BatchTimeout: 50 * time.Millisecond,
		RequiredAcks: kafka.RequireAll,
		Async:        true,
	}
	mqProducer, err := kafkago.NewProducer(conf)
	if err != nil {
		logger.Log().Err(err).Msg("failed to create kafka producer")
		return
	}

	route := newRouterSrv(mqProducer)
	conv := newConversationSrv()
	msg := newMessageSrv(conv, route)

	svc = &ServiceContext{
		Route: route,
		Conv:  conv,
		Msg:   msg,
	}
}

func Hints() *ServiceContext {
	return svc
}
