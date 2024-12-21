package net

type ServerContext struct {
	online int32

	processor   *Processor
	ConnManager *ClientManager
	Retry       *RetryDelivery
}

func InitServerContext() *ServerContext {
	svc := &ServerContext{
		ConnManager: newClientManager(),
	}

	svc.processor = newProcessor(svc)

	svc.Retry = newRetryDelivery(svc)
	// todo retry启动
	return svc
}
