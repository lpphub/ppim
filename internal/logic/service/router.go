package service

import (
	"context"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/segmentio/kafka-go"
	"ppim/internal/chatlib"
	"ppim/internal/logic/global"
	"ppim/internal/logic/types"
	"ppim/pkg/kafka/producer"
	"strings"
)

type RouterSrv struct {
	mq *producer.Producer
}

const (
	cacheRouteUid = "online:%s"
)

func newRouterSrv() *RouterSrv {
	mqProducer, _ := producer.NewProducer(producer.WithBrokers(global.Conf.Kafka.Brokers))
	return &RouterSrv{
		mq: mqProducer,
	}
}

func (s *RouterSrv) Register(ctx context.Context, ol *types.RouteDTO) error {
	err := global.Redis.SAdd(ctx, s.genRouteKey(ol.Uid), s.buildVal(ol)).Err()
	return err
}

func (s *RouterSrv) UnRegister(ctx context.Context, ol *types.RouteDTO) error {
	err := global.Redis.SAdd(ctx, s.genRouteKey(ol.Uid), s.buildVal(ol)).Err()
	return err
}

func (s *RouterSrv) RouteChat(ctx context.Context, routeKeys []string, msg *types.MessageDTO) error {
	messageSlice := make([]kafka.Message, 0, len(routeKeys))
	for _, key := range routeKeys {
		route := strings.Split(key, "_")

		bytes, _ := jsoniter.Marshal(msg)
		dd := &chatlib.DeliverMsg{
			CMD:     "chat",
			ToUID:   []string{route[0]},
			Content: bytes,
		}
		message := kafka.Message{
			Topic: route[2],
			Value: dd.ToJsonBytes(),
		}
		messageSlice = append(messageSlice, message)
	}
	return s.mq.SendMessage(ctx, messageSlice...)
}

func (s *RouterSrv) genRouteKey(uid string) string {
	return fmt.Sprintf(cacheRouteUid, uid)
}

func (s *RouterSrv) buildVal(ol *types.RouteDTO) string {
	return fmt.Sprintf("%s_%s_%s", ol.Uid, ol.Did, ol.Topic)
}
