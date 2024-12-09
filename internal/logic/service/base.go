package service

import (
	"ppim/internal/logic/global"
	"ppim/pkg/ext"
	"ppim/pkg/kafka/producer"
)

type ServiceContext struct {
	MsgSrv    *MessageSrv
	RouterSrv *RouterSrv
	ConvSrv   *ConversationSrv
}

var svc *ServiceContext

func init() {
	mqp, _ := producer.NewProducer(producer.WithBrokers(global.Conf.Kafka.Brokers))

	svc = &ServiceContext{
		MsgSrv: &MessageSrv{},
		RouterSrv: &RouterSrv{
			mq: mqp,
		},
		ConvSrv: &ConversationSrv{
			segmentLock: *ext.NewSegmentLock(20),
		},
	}
}

func Inst() *ServiceContext {
	return svc
}
