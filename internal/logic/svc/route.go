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
	"ppim/pkg/util"
	"strings"
)

type RouteSrv struct {
	mq    *kafkago.Producer
	cache *redis.Client
}

const (
	// 缓存路由：route:123 pc_d003 topic1 ios_i201 topic2
	cacheRouteUid = "route:%s"
)

func newRouterSrv(mq *kafkago.Producer) *RouteSrv {
	return &RouteSrv{
		mq:    mq,
		cache: global.Redis,
	}
}

func (s *RouteSrv) Online(ctx context.Context, ol *types.RouteDTO) error {
	err := s.cache.HSet(ctx, s.genRouteKey(ol.Uid), ol.Did, ol.Topic).Err()
	if err != nil {
		logger.Err(ctx, err, "route online error")
	}
	return err
}

func (s *RouteSrv) Offline(ctx context.Context, ol *types.RouteDTO) error {
	err := s.cache.HDel(ctx, s.genRouteKey(ol.Uid), ol.Did).Err()
	if err != nil {
		logger.Err(ctx, err, "route offline error")
	}
	return err
}

func (s *RouteSrv) batchGetOnline(ctx context.Context, uids []string) ([]*redis.MapStringStringCmd, error) {
	pipe := s.cache.Pipeline()
	cmds := make([]*redis.MapStringStringCmd, len(uids))
	for i, uid := range uids {
		cmds[i] = pipe.HGetAll(ctx, s.genRouteKey(uid))
	}
	_, err := pipe.Exec(ctx)
	return cmds, err
}

func (s *RouteSrv) RouteChat(ctx context.Context, msg *types.MessageDTO, receivers []string) error {
	// 在线用户设备路由, 离线用户UID
	var onlineSlice, offlineSlice []string
	for _, chunk := range util.Partition(receivers, 300) {
		cmds, err := s.batchGetOnline(ctx, chunk)
		if err != nil {
			logger.Err(ctx, err, fmt.Sprintf("batch get online error: %v", chunk))
			continue
		}
		for i, uid := range chunk {
			if online, _ := cmds[i].Result(); len(online) > 0 {
				for did, topic := range online {
					if did != msg.FromDID || uid != msg.FromUID { // 排除发送者同一设备，而不同设备时则接收消息
						onlineSlice = append(onlineSlice, fmt.Sprintf("%s#%s", uid, topic))
					}
				}
			} else {
				offlineSlice = append(offlineSlice, uid)
			}
		}
	}

	// 在线投递(去重合并同一消息同一topic)
	if len(onlineSlice) > 0 {
		onlineSlice = util.RemoveDup(onlineSlice)
		topicReceivers := make(map[string][]string)
		for _, key := range onlineSlice {
			route := strings.Split(key, "#")
			topicReceivers[route[1]] = append(topicReceivers[route[1]], route[0])
		}

		messageSlice := make([]kafka.Message, 0, len(topicReceivers))
		for t, us := range topicReceivers {
			message := kafka.Message{
				Topic: t,
				Value: (&chatlib.DeliveryMsg{
					CMD:       chatlib.DeliveryChat,
					FromUID:   msg.FromUID,
					Receivers: us,
					Chat:      s.convert(msg),
				}).ToJsonBytes(),
			}
			messageSlice = append(messageSlice, message)
		}

		if err := s.mq.SendMessage(ctx, messageSlice...); err != nil {
			return err
		}
	}

	if len(offlineSlice) > 0 {
		// todo 消息离线通知
		logger.Warnf(ctx, "offline push: %v", offlineSlice)
	}
	return nil
}

func (s *RouteSrv) convert(msg *types.MessageDTO) *chatlib.ChatMsg {
	var chat chatlib.ChatMsg
	_ = copier.Copy(&chat, msg)
	return &chat
}

func (s *RouteSrv) genRouteKey(uid string) string {
	return fmt.Sprintf(cacheRouteUid, uid)
}
