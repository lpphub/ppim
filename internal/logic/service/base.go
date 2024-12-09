package service

type ServiceContext struct {
	MsgSrv    *MessageSrv
	RouterSrv *RouterSrv
}

var svc *ServiceContext

func init() {
	svc = &ServiceContext{
		MsgSrv:    newMessageSrv(),
		RouterSrv: newRouterSrv(),
	}
}

func Inst() *ServiceContext {
	return svc
}
