package net

type ServerContext struct {
	online int32

	processor    *Processor
	ConnManager  *ClientManager
	RetryManager *RetryManager
}

func InitServerContext() *ServerContext {
	svc := &ServerContext{
		ConnManager: newClientManager(),
	}

	svc.processor = newProcessor(svc)
	svc.RetryManager = newRetryManager(svc, 2)
	return svc
}

func (svc *ServerContext) StartBackground() {
	svc.RetryManager.Start()
}
