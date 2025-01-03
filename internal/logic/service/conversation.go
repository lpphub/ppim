package service

import (
	"context"
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/lpphub/golib/gowork"
	"github.com/lpphub/golib/logger"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
	"go.mongodb.org/mongo-driver/mongo"
	"ppim/internal/chatlib"
	"ppim/internal/logic/global"
	"ppim/internal/logic/store"
	"ppim/internal/logic/types"
	"ppim/pkg/ext"
	"time"
)

// ConversationSrv
/**
 * 写扩散：每个用户对应一个timeline, 消息到达后每个接收者更新自身timeline
 * 读扩散：一个会话对应一个timeline，消息到达后更新此会话最新timeline
 */
type ConversationSrv struct {
	works       *gowork.Pool
	segmentLock *ext.SegmentRWLock
}

func newConversationSrv() *ConversationSrv {
	return &ConversationSrv{
		works:       gowork.NewPool(100),
		segmentLock: ext.NewSegmentLock(20),
	}
}

/**
 * 1.会话最新消息
 * 2.用户最近会话列表：只保留100条
 * 3.用户会话详细：未读消息数，最新消息, 置顶，免打扰，已读消息
 *
 * conv:recent:msg:{convID} -> string -> msg
 * conv:recent:{uid} -> sortedset -> convID,sendTime
 * conv:recent:{uid}:{convID} -> hash -> unreadCount,pin,mute,readMsgId
 */
const (
	CacheConvRecentMsg  = "conv:recent:msg:%s"
	CacheConvRecent     = "conv:recent:%s"
	CacheConvRecentInfo = "conv:recent:%s:%s"

	CacheFieldConvUnreadCount = "unreadCount"
	CacheFieldConvLastMsgId   = "lastMsgId"
	CacheFieldConvPin         = "pin"
	CacheFieldConvMute        = "mute"
	CacheFieldConvReadMsgId   = "readMsgId"
)

func (c *ConversationSrv) IndexRecent(ctx context.Context, msg *types.MessageDTO, receivers []string) error {
	msgJson, _ := jsoniter.MarshalToString(msg)
	global.Redis.Set(ctx, fmt.Sprintf(CacheConvRecentMsg, msg.ConversationID), msgJson, 30*24*time.Hour)

	for _, uid := range receivers {
		_ = c.works.Submit(func() {
			c.indexWithLock(ctx, msg, uid)
		})
	}
	return nil
}

func (c *ConversationSrv) indexWithLock(ctx context.Context, msg *types.MessageDTO, uid string) {
	// todo 集群模式下，分布式锁
	index := cast.ToInt(chatlib.DigitizeUID(uid))
	c.segmentLock.Lock(index)
	defer c.segmentLock.Unlock(index)

	if err := c.cacheStoreRecent(ctx, uid, msg); err != nil {
		logger.Err(ctx, err, fmt.Sprintf("conversation cache recent: uid=%s", uid))
		return
	}

	// todo 优化：异步周期性从redis持久化至存储层，减少数据存储操作
	conv, err := new(store.Conversation).GetOne(ctx, uid, msg.ConversationID)
	if err != nil && errors.Is(err, mongo.ErrNoDocuments) {
		conv = &store.Conversation{
			ConversationID:   msg.ConversationID,
			ConversationType: msg.ConversationType,
			UID:              uid,
			UnreadCount:      1,
			LastMsgId:        msg.MsgID,
			LastMsgSeq:       msg.MsgSeq,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		}
		_ = conv.Insert(ctx)
	} else {
		if msg.FromUID != uid {
			conv.UnreadCount++
		}
		if conv.LastMsgSeq < msg.MsgSeq {
			conv.LastMsgId = msg.MsgID
			conv.LastMsgSeq = msg.MsgSeq
			conv.FromUID = msg.FromUID
		}
		conv.UpdatedAt = time.Now()
		_ = conv.Update(ctx)
	}
}

// 缓存用户最近会话
func (c *ConversationSrv) cacheStoreRecent(ctx context.Context, uid string, msg *types.MessageDTO) error {
	cacheKey := fmt.Sprintf(CacheConvRecent, uid)
	cacheInfoKey := fmt.Sprintf(CacheConvRecentInfo, uid, msg.ConversationID)

	pipe := global.Redis.Pipeline()
	// 用户最新会话
	pipe.ZAdd(ctx, cacheKey, redis.Z{Score: float64(msg.SendTime), Member: msg.ConversationID})

	// 会话最新消息
	pipe.HSet(ctx, cacheInfoKey, CacheFieldConvLastMsgId, msg.MsgID)

	if uid != msg.FromUID {
		// 未读消息数
		pipe.HIncrBy(ctx, cacheInfoKey, CacheFieldConvUnreadCount, 1)
	}

	// 会话数量超过100，删除最早一条
	count, _ := pipe.ZCard(ctx, cacheKey).Result()
	if count > 100 {
		pipe.ZRemRangeByRank(ctx, cacheKey, 0, 0)
		pipe.HDel(ctx, cacheInfoKey)
	}
	_, err := pipe.Exec(ctx)
	return err
}

func (c *ConversationSrv) CacheQueryRecent(ctx context.Context, uid string) ([]*types.RecentConvVO, error) {
	ids, err := global.Redis.ZRevRange(ctx, fmt.Sprintf(CacheConvRecent, uid), 0, 200).Result()
	if err != nil {
		return nil, err
	}
	logger.Infof(ctx, "recent conv ids=%v", ids)

	var list []*types.RecentConvVO
	for _, id := range ids {
		pipe := global.Redis.Pipeline()
		pipe.Get(ctx, fmt.Sprintf(CacheConvRecentMsg, id))
		pipe.HMGet(ctx, fmt.Sprintf(CacheConvRecentInfo, uid, id), CacheFieldConvUnreadCount, CacheFieldConvPin, CacheFieldConvMute)
		cmds, err := pipe.Exec(ctx)
		if err != nil {
			logger.Err(ctx, err, "conv: cache query recent")
			continue
		}

		vo := &types.RecentConvVO{}
		recentMsg, _ := cmds[0].(*redis.StringCmd).Result()
		if recentMsg != "" {
			var mt types.MessageDTO
			_ = jsoniter.UnmarshalFromString(recentMsg, &mt)

			vo.ConversationID = mt.ConversationID
			vo.ConversationType = mt.ConversationType
			vo.FromUID = mt.FromUID
			vo.LastMsgID = mt.MsgID
			vo.LastMsgSeq = mt.MsgSeq
			vo.Version = mt.CreatedAt
			vo.LastMsg = &mt
		}

		fields, _ := cmds[1].(*redis.SliceCmd).Result()
		logger.Infof(ctx, "conv: cache query recent fields=%v", fields)
		if len(fields) == 3 {
			if fields[0] != nil {
				vo.UnreadCount = cast.ToUint64(fields[0])
			}
			if fields[1] != nil {
				vo.Pin = fields[1].(bool)
			}
			if fields[2] != nil {
				vo.Mute = fields[2].(bool)
			}
		}
		list = append(list, vo)
	}
	return list, nil
}
