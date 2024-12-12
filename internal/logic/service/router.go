package service

import (
	"context"
	"fmt"
	"github.com/lpphub/golib/logger"
	"github.com/segmentio/kafka-go"
	"ppim/internal/chatlib"
	"ppim/internal/logic/global"
	"ppim/internal/logic/types"
	"ppim/pkg/kafka/producer"
	"strings"
	"time"
)

type RouterSrv struct {
	mq *producer.Producer
}

const (
	cacheRouteUid = "route:%s"
)

func newRouterSrv(mq *producer.Producer) *RouterSrv {
	return &RouterSrv{
		mq: mq,
	}
}

func (s *RouterSrv) Register(ctx context.Context, ol *types.RouteDTO) error {
	err := global.Redis.SAdd(ctx, s.genRouteKey(ol.Uid), s.buildVal(ol)).Err()
	return err
}

func (s *RouterSrv) UnRegister(ctx context.Context, ol *types.RouteDTO) error {
	err := global.Redis.SRem(ctx, s.genRouteKey(ol.Uid), s.buildVal(ol)).Err()
	return err
}

func (s *RouterSrv) RouteDeliver(ctx context.Context, routeKeys []string, msg *types.MessageDTO) error {
	messageSlice := make([]kafka.Message, 0, len(routeKeys))
	for _, key := range routeKeys {
		route := strings.Split(key, "_")

		dd := &chatlib.DeliverMsg{
			CMD:     chatlib.DeliverChat,
			ToUID:   route[0],
			ChatMsg: msg,
		}
		message := kafka.Message{
			Topic: route[2],
			Value: dd.ToJsonBytes(),
		}
		messageSlice = append(messageSlice, message)
	}

	t := time.Now()
	err := s.mq.SendMessage(ctx, messageSlice...)
	logger.Infof(ctx, "MQ投递耗时: %d", time.Since(t).Milliseconds())

	return err
}

func (s *RouterSrv) genRouteKey(uid string) string {
	return fmt.Sprintf(cacheRouteUid, uid)
}

func (s *RouterSrv) buildVal(ol *types.RouteDTO) string {
	return fmt.Sprintf("%s_%s_%s", ol.Uid, ol.Did, ol.Topic)
}
