package service

import (
	"github.com/lpphub/golib/logger"
	"ppim/internal/logic/global"
	"ppim/pkg/kafka/producer"
)

type ServiceContext struct {
	MsgSrv    *MessageSrv
	RouterSrv *RouterSrv
}

var svc *ServiceContext

func LoadService() {
	mqProducer, err := producer.NewProducer(producer.WithBrokers(global.Conf.Kafka.Brokers))
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
