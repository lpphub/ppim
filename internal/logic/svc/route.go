package svc

import (
	"context"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/lpphub/golib/logger"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
	"ppim/internal/chatlib"
	"ppim/internal/logic/global"
	"ppim/internal/logic/types"
	"ppim/pkg/kafkago"
	"strings"
	"time"
)

type RouteSrv struct {
	mq         *kafkago.Producer
	cacheStore *redis.Client
}

const (
	// 缓存路由：route:123 pc_d003 topic1 ios_i201 topic2
	cacheRouteUid = "route:%s"
)

func newRouterSrv(mq *kafkago.Producer) *RouteSrv {
	return &RouteSrv{
		mq:         mq,
		cacheStore: global.Redis,
	}
}

func (s *RouteSrv) Online(ctx context.Context, ol *types.RouteDTO) error {
	err := s.cacheStore.HSet(ctx, s.genRouteKey(ol.Uid), ol.Did, ol.Topic).Err()
	if err != nil {
		logger.Err(ctx, err, "route online error")
	}
	return err
}

func (s *RouteSrv) Offline(ctx context.Context, ol *types.RouteDTO) error {
	err := s.cacheStore.HDel(ctx, s.genRouteKey(ol.Uid), ol.Did).Err()
	if err != nil {
		logger.Err(ctx, err, "route offline error")
	}
	return err
}

func (s *RouteSrv) BatchGetOnline(ctx context.Context, uids []string) ([]*redis.MapStringStringCmd, error) {
	pipe := s.cacheStore.Pipeline()
	cmd := make([]*redis.MapStringStringCmd, len(uids))
	for i, uid := range uids {
		cmd[i] = pipe.HGetAll(ctx, s.genRouteKey(uid))
	}
	_, err := pipe.Exec(ctx)
	if err != nil {
		return nil, err
	}
	return cmd, nil
}

func (s *RouteSrv) RouteChat(ctx context.Context, routeKeys []string, msg *types.MessageDTO) error {
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
		dm := &chatlib.DeliveryMsg{
			CMD:       chatlib.DeliveryChat,
			FromUID:   msg.FromUID,
			Receivers: us,
			Chat:      s.convert(msg),
		}
		message := kafka.Message{
			Value: dm.ToJsonBytes(),
			Topic: t,
		}
		messageSlice = append(messageSlice, message)
	}

	t := time.Now()
	err := s.mq.SendMessage(ctx, messageSlice...)
	logger.Infof(ctx, "MQ produce cost_ms: %d", time.Since(t).Milliseconds())
	return err
}

func (s *RouteSrv) convert(msg *types.MessageDTO) *chatlib.ChatMsg {
	var chat chatlib.ChatMsg
	_ = copier.Copy(&chat, msg)
	return &chat
}

func (s *RouteSrv) genRouteKey(uid string) string {
	return fmt.Sprintf(cacheRouteUid, uid)
}
