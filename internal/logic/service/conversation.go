package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/lpphub/golib/gowork"
	"github.com/lpphub/golib/logger"
	"github.com/redis/go-redis/v9"
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
	segmentLock ext.SegmentLock
}

func newConversationSrv() *ConversationSrv {
	return &ConversationSrv{
		works:       gowork.NewPool(100),
		segmentLock: *ext.NewSegmentLock(20),
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

func (c *ConversationSrv) IndexRecent(ctx context.Context, msg *types.MessageDTO, uidSlice []string) error {
	// 发送者 + 接收者
	uidSlice = append(uidSlice, msg.FromID)
	for _, uid := range uidSlice {
		_ = c.works.Submit(func() {
			c.indexWithLock(ctx, msg, uid)
		})
	}
	return nil
}

func (c *ConversationSrv) indexWithLock(ctx context.Context, msg *types.MessageDTO, uid string) {
	// todo 集群模式下，分布式锁
	index := chatlib.DigitizeUID(uid)
	c.segmentLock.Lock(index)
	defer c.segmentLock.Unlock(index)

	if err := c.cacheRecent(ctx, uid, msg); err != nil {
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
		if msg.FromID != uid {
			conv.UnreadCount++
		}
		if conv.LastMsgSeq < msg.MsgSeq {
			conv.LastMsgId = msg.MsgID
			conv.LastMsgSeq = msg.MsgSeq
			conv.FromID = msg.FromID
		}
		conv.UpdatedAt = time.Now()
		_ = conv.Update(ctx)
	}
}

// 缓存用户最近会话
func (c *ConversationSrv) cacheRecent(ctx context.Context, uid string, msg *types.MessageDTO) error {
	cacheKey := fmt.Sprintf(CacheConvRecent, uid)
	cacheInfoKey := fmt.Sprintf(CacheConvRecentInfo, uid, msg.ConversationID)

	pipe := global.Redis.Pipeline()
	// 用户最新会话
	pipe.ZAdd(ctx, cacheKey, redis.Z{Score: float64(msg.SendTime), Member: msg.ConversationID})

	// 会话最新消息
	pipe.HSet(ctx, cacheInfoKey, CacheFieldConvLastMsgId, msg.MsgID)

	if uid != msg.FromID {
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
