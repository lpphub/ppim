package net

type ServerContext struct {
	online int32

	ConnManager *ClientManager
	processor   *Processor
}

func InitServerContext() *ServerContext {
	svc := &ServerContext{
		ConnManager: newClientManager(),
	}

	svc.processor = newProcessor(svc)
	return svc
}
