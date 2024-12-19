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

type RouteSrv struct {
	mq *producer.Producer
}

const (
	// 缓存路由：route:123 pc_d003 topic1 ios_i201 topic2
	cacheRouteUid = "route:%s"
)

func newRouterSrv(mq *producer.Producer) *RouteSrv {
	return &RouteSrv{
		mq: mq,
	}
}

func (s *RouteSrv) Online(ctx context.Context, ol *types.RouteDTO) error {
	err := global.Redis.HSet(ctx, s.genRouteKey(ol.Uid), ol.Did, ol.Topic).Err()
	if err != nil {
		logger.Err(ctx, err, "route online error")
	}
	return err
}

func (s *RouteSrv) Offline(ctx context.Context, ol *types.RouteDTO) error {
	err := global.Redis.HDel(ctx, s.genRouteKey(ol.Uid), ol.Did).Err()
	if err != nil {
		logger.Err(ctx, err, "route offline error")
	}
	return err
}

func (s *RouteSrv) RouteDelivery(ctx context.Context, routeKeys []string, msg *types.MessageDTO) error {
	topicReceivers := make(map[string][]string)
	for _, key := range routeKeys {
		route := strings.Split(key, "#") // key = uid#topic
		if r, ok := topicReceivers[route[1]]; ok {
			topicReceivers[route[1]] = append(r, route[0])
		} else {
			topicReceivers[route[1]] = []string{route[0]}
		}
	}

	messageSlice := make([]kafka.Message, 0, len(routeKeys))
	for t, us := range topicReceivers {
		dd := &chatlib.DeliveryMsg{
			CMD:       chatlib.DeliveryChat,
			Receivers: us,
			ChatMsg:   msg,
		}
		message := kafka.Message{
			Value: dd.ToJsonBytes(),
			Topic: t,
		}
		messageSlice = append(messageSlice, message)
	}

	t := time.Now()
	err := s.mq.SendMessage(ctx, messageSlice...)
	logger.Infof(ctx, "MQ produce cost_ms: %d", time.Since(t).Milliseconds())
	return err
}

func (s *RouteSrv) genRouteKey(uid string) string {
	return fmt.Sprintf(cacheRouteUid, uid)
}

func (s *RouteSrv) buildVal(ol *types.RouteDTO) string {
	return fmt.Sprintf("%s_%s_%s", ol.Uid, ol.Did, ol.Topic)
}
