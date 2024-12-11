package net

type ServerContext struct {
	online int32

	ConnManager *ClientManager
	processor   *Processor
}

func NewServerContext() *ServerContext {
	svc := &ServerContext{
		ConnManager: newClientManager(),
	}

	svc.processor = newProcessor(svc)
	return svc
}
