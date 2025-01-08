package svc

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
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
	workers     *gowork.Pool
	segmentLock *ext.SegmentRWLock
}

func newConversationSrv() *ConversationSrv {
	return &ConversationSrv{
		workers:     gowork.NewPool(100),
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
 * conv:recent:{uid}:{convID} -> hash -> unreadCount,pin,mute,readMsgSeq
 */
const (
	CacheConvRecentMsg  = "conv:recent:msg:%s"
	CacheConvRecent     = "conv:recent:%s"
	CacheConvRecentInfo = "conv:recent:%s:%s"

	ConvFieldUnreadCount = "unreadCount"
	ConvFieldPin         = "pin"
	ConvFieldMute        = "mute"
	ConvFieldLastMsgSeq  = "lastMsgSeq"
	ConvFieldReadMsgSeq  = "readMsgSeq"

	recentConvMaxSize = 500
)

func (c *ConversationSrv) IndexRecent(ctx context.Context, msg *types.MessageDTO, receivers []string) error {
	msgJson, _ := jsoniter.MarshalToString(msg)
	global.Redis.Set(ctx, fmt.Sprintf(CacheConvRecentMsg, msg.ConversationID), msgJson, 30*24*time.Hour)

	for _, uid := range receivers {
		_ = c.workers.Submit(func() {
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
	cacheInfoKey := c.getConvCacheKey(uid, msg.ConversationID)

	rdb := global.Redis
	pipe := rdb.Pipeline()
	// 会话最新消息
	pipe.HSet(ctx, cacheInfoKey, ConvFieldLastMsgSeq, msg.MsgSeq)
	if uid != msg.FromUID {
		// 未读消息数
		pipe.HIncrBy(ctx, cacheInfoKey, ConvFieldUnreadCount, 1)
	}
	// 用户最新会话
	pipe.ZAdd(ctx, cacheKey, redis.Z{Score: float64(msg.SendTime), Member: msg.ConversationID})
	_, err := pipe.Exec(ctx)
	if err != nil {
		return err
	}

	// 获取会话数量，判断是否超过限制
	count, _ := rdb.ZCard(ctx, cacheKey).Result()
	if count > recentConvMaxSize {
		first, _ := rdb.ZRangeWithScores(ctx, cacheKey, 0, 0).Result()
		if len(first) > 0 {
			rdb.ZRem(ctx, cacheInfoKey, first[0].Member)
			rdb.HDel(ctx, c.getConvCacheKey(uid, first[0].Member.(string)))
		}
	}
	return nil
}

func (c *ConversationSrv) cacheQueryRecent(ctx context.Context, uid string) ([]*types.ConvRecentDTO, error) {
	cids, err := global.Redis.ZRevRange(ctx, fmt.Sprintf(CacheConvRecent, uid), 0, recentConvMaxSize).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get recent conversation IDs: %v", err)
	}
	logger.Infof(ctx, "recent conv ids=%v", cids)
	if len(cids) == 0 {
		return nil, redis.Nil
	}

	pipe := global.Redis.Pipeline()
	type cmdPair struct {
		msgCmd  *redis.StringCmd
		infoCmd *redis.SliceCmd
	}
	cmds := make([]cmdPair, len(cids))
	// 批量添加命令到 Pipeline
	for i, cid := range cids {
		cmds[i].msgCmd = pipe.Get(ctx, fmt.Sprintf(CacheConvRecentMsg, cid))
		cmds[i].infoCmd = pipe.HMGet(ctx, c.getConvCacheKey(uid, cid), ConvFieldUnreadCount, ConvFieldPin, ConvFieldMute)
	}
	_, err = pipe.Exec(ctx)
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, fmt.Errorf("failed to execute pipeline: %v", err)
	}

	var list []*types.ConvRecentDTO
	for i := range cids {
		conv := new(types.ConvRecentDTO)

		// 会话最近消息
		recentMsg, _ := cmds[i].msgCmd.Result()
		if recentMsg != "" {
			var mt types.MessageDTO
			_ = jsoniter.UnmarshalFromString(recentMsg, &mt)

			conv.ConversationID = mt.ConversationID
			conv.ConversationType = mt.ConversationType
			conv.LastMsg = &mt
			conv.LastMsgID = mt.MsgID
			conv.LastMsgSeq = mt.MsgSeq
			conv.Version = mt.CreatedAt
		}
		// 会话详情信息
		fields, _ := cmds[i].infoCmd.Result()
		if len(fields) == 3 {
			if fields[0] != nil {
				conv.UnreadCount = cast.ToUint64(fields[0])
			}
			if fields[1] != nil {
				conv.Pin = fields[1].(bool)
			}
			if fields[2] != nil {
				conv.Mute = fields[2].(bool)
			}
		}
		list = append(list, conv)
	}
	return list, nil
}

func (c *ConversationSrv) getConvCacheKey(uid, convID string) string {
	return fmt.Sprintf(CacheConvRecentInfo, uid, convID)
}

func (c *ConversationSrv) GetRecentByUID(ctx *gin.Context, uid string) ([]*types.ConvRecentDTO, error) {
	if list, err := c.cacheQueryRecent(ctx, uid); err == nil {
		return list, nil
	} else {
		logger.Err(ctx, err, fmt.Sprintf("conversation recent cache query: uid=%s", uid))
	}

	// todo 可以全用缓存，不查db
	data, err := new(store.Conversation).ListRecent(ctx, uid, recentConvMaxSize)
	if err != nil {
		return nil, err
	}

	list := make([]*types.ConvRecentDTO, 0, len(data))
	msgIds := make([]string, 0, len(data))
	for _, d := range data {
		list = append(list, &types.ConvRecentDTO{
			ConversationID:   d.ConversationID,
			ConversationType: d.ConversationType,
			UnreadCount:      d.UnreadCount,
			Mute:             d.Mute,
			Pin:              d.Pin,
			LastMsgID:        d.LastMsgId,
			Version:          d.CreatedAt.UnixMilli(),
		})
		msgIds = append(msgIds, d.LastMsgId)
	}

	msgList, err := new(store.Message).ListByMsgIds(ctx, msgIds)
	if err != nil {
		return nil, err
	}
	msgMap := make(map[string]*types.MessageDTO, len(msgList))
	for _, m := range msgList {
		var md types.MessageDTO
		_ = copier.Copy(&md, m)
		md.CreatedAt = m.CreatedAt.UnixMilli()
		md.SendTime = m.SendTime.UnixMilli()
		msgMap[m.MsgID] = &md
	}
	for _, cv := range list {
		if md, ok := msgMap[cv.LastMsgID]; ok {
			cv.LastMsg = md
			cv.LastMsgSeq = md.MsgSeq
			cv.Version = md.CreatedAt
		}
	}
	return list, nil
}

func (c *ConversationSrv) SetAttribute(ctx context.Context, attr types.ConvAttributeDTO) (err error) {
	cacheKey := c.getConvCacheKey(attr.UID, attr.ConversationID)
	switch attr.Attribute {
	case ConvFieldPin:
		global.Redis.HSet(ctx, cacheKey, ConvFieldPin, attr.Pin)
		err = new(store.Conversation).UpdatePin(ctx, attr.UID, attr.ConversationID, attr.Pin)
	case ConvFieldMute:
		global.Redis.HSet(ctx, cacheKey, ConvFieldMute, attr.Mute)
		err = new(store.Conversation).UpdateMute(ctx, attr.UID, attr.ConversationID, attr.Mute)
	case ConvFieldUnreadCount:
		global.Redis.HSet(ctx, cacheKey, ConvFieldUnreadCount, attr.UnreadCount)
		err = new(store.Conversation).UpdateUnreadCount(ctx, attr.UID, attr.ConversationID, attr.UnreadCount)
	default:
		return errors.New("invalid op type")
	}
	if err != nil {
		return err
	}
	global.Redis.ZAdd(ctx, fmt.Sprintf(CacheConvRecent, attr.UID), redis.Z{Score: float64(time.Now().UnixMilli()), Member: attr.ConversationID})
	return
}
