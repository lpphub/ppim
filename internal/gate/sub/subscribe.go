package sub

import (
	"context"
	"github.com/lpphub/golib/logger"
	"ppim/internal/clib"
	"ppim/internal/gate/net"
)

type Subscriber struct {
	svc *net.ServerContext
}

func subscriber(ctx context.Context, topic string) {

}

func (s *Subscriber) deliver(ctx context.Context, msg clib.DeliverMsg) {
	for _, uid := range msg.ToUID {
		clients := s.svc.ConnManager.GetWithUID(uid)

		for _, client := range clients {
			_, err := client.Write(msg.MsgData)
			if err != nil {
				logger.Err(ctx, err, "write to client error")
			}
		}
	}

}
