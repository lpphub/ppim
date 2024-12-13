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
	// 缓存格式：route:123 pc_d003 topic1 ios_i201 topic2
	cacheRouteUid = "route:%s"
)

func newRouterSrv(mq *producer.Producer) *RouterSrv {
	return &RouterSrv{
		mq: mq,
	}
}

func (s *RouterSrv) Online(ctx context.Context, ol *types.RouteDTO) error {
	err := global.Redis.HSet(ctx, s.genRouteKey(ol.Uid), ol.Did, ol.Topic).Err()
	return err
}

func (s *RouterSrv) Offline(ctx context.Context, ol *types.RouteDTO) error {
	err := global.Redis.HDel(ctx, s.genRouteKey(ol.Uid), ol.Did).Err()
	return err
}

func (s *RouterSrv) RouteDeliver(ctx context.Context, routeKeys []string, msg *types.MessageDTO) error {
	messageSlice := make([]kafka.Message, 0, len(routeKeys))
	for _, key := range routeKeys {
		route := strings.Split(key, "#")

		dd := &chatlib.DeliverMsg{
			CMD:     chatlib.DeliverChat,
			ToUID:   route[0],
			ChatMsg: msg,
		}
		message := kafka.Message{
			Topic: route[1],
			Value: dd.ToJsonBytes(),
		}
		messageSlice = append(messageSlice, message)
	}

	t := time.Now()
	err := s.mq.SendMessage(ctx, messageSlice...)
	logger.Infof(ctx, "MQ produce cost_ms: %d", time.Since(t).Milliseconds())
	return err
}

func (s *RouterSrv) genRouteKey(uid string) string {
	return fmt.Sprintf(cacheRouteUid, uid)
}

func (s *RouterSrv) buildVal(ol *types.RouteDTO) string {
	return fmt.Sprintf("%s_%s_%s", ol.Uid, ol.Did, ol.Topic)
}
