package service

import (
	"github.com/lpphub/golib/logger"
	"ppim/internal/logic/global"
	"ppim/pkg/kafka/producer"
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
	mqProducer, err := producer.NewProducer(producer.WithBrokers(global.Conf.Kafka.Brokers),
		producer.WithBatchTimeout(50*time.Millisecond), producer.WithAsync(true))
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

func Hint() *ServiceContext {
	return svc
}
