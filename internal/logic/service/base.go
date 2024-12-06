package service

import (
	"ppim/pkg/ext"
)

type ServiceContext struct {
	MsgSrv    *MessageSrv
	OnlineSrv *OnlineSrv
	ConvSrv   *ConversationSrv
}

var svc *ServiceContext

func init() {
	svc = &ServiceContext{
		MsgSrv:    &MessageSrv{},
		OnlineSrv: &OnlineSrv{},
		ConvSrv: &ConversationSrv{
			segmentLock: *ext.NewSegmentLock(20),
		},
	}
}

func Inst() *ServiceContext {
	return svc
}
