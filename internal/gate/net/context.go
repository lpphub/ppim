package net

type ServerContext struct {
	online int32

	connManager *ClientManager
	processor   *Processor
}

func NewServerContext() *ServerContext {
	svc := &ServerContext{
		connManager: newClientManager(),
	}

	svc.processor = newProcessor(svc)
	return svc
}
