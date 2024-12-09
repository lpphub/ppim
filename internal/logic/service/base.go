package service

import (
	"github.com/lpphub/golib/logger"
	"ppim/internal/logic/global"
	"ppim/pkg/kafka/producer"
	"time"
)

type ServiceContext struct {
	MsgSrv    *MessageSrv
	RouterSrv *RouterSrv
}

var svc *ServiceContext

func LoadService() {
	// kafka flush msg every 100ms
	mqProducer, err := producer.NewProducer(producer.WithBrokers(global.Conf.Kafka.Brokers),
		producer.WithBatchTimeout(10*time.Millisecond), producer.WithAsync(false))
	if err != nil {
		logger.Log().Err(err).Msg("failed to create kafka producer")
		return
	}

	route := newRouterSrv(mqProducer)

	svc = &ServiceContext{
		MsgSrv:    newMessageSrv(route),
		RouterSrv: route,
	}
}

func Inst() *ServiceContext {
	return svc
}
