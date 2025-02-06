package svc

import (
	"context"
	"errors"
	"fmt"
	"github.com/jinzhu/copier"
	jsoniter "github.com/json-iterator/go"
	"github.com/lpphub/golib/logger"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"ppim/internal/logic/global"
	"ppim/internal/logic/store"
	"ppim/internal/logic/types"
	"ppim/pkg/ext"
	"ppim/pkg/util"
	"time"
)

// ConversationSrv
/**
 * 写扩散：每个用户对应一个timeline, 消息到达后每个接收者更新自身timeline
 * 读扩散：一个会话对应一个timeline，消息到达后更新此会话最新timeline
 */
type ConversationSrv struct {
	cache      *redis.Client
	batchStore *ext.BatchProcessor[*convStoreData]
}

type convStoreData struct {
	UID     string
	LastMsg *types.MessageDTO
	Count   int
}

/**
 * 1.会话最新消息
 * 2.用户会话列表：可按时间或数量限制
 * 3.用户会话详细：未读消息数，最新消息, 置顶，免打扰，已读消息
 *
 * conv:last_msg:{convID} -> string -> msg
 * conv:list:{uid} -> sortedset -> convID,sendTime
 * conv:info:{uid}:{convID} -> hash -> unreadCount,pin,mute,readMsgSeq
 */
const (
	CacheConvLastMsg = "conv:last_msg:%s"
	CacheConvList    = "conv:list:%s"
	CacheConvInfo    = "conv:info:%s:%s"

	ConvFieldCreatedAt   = "createdAt"
	ConvFieldUnreadCount = "unreadCount"
	ConvFieldPin         = "pin"
	ConvFieldMute        = "mute"
	ConvFieldLastMsgSeq  = "lastMsgSeq"
	ConvFieldDeleted     = "deleted"
	ConvFieldReadMsgSeq  = "readMsgSeq"
)

func newConversationSrv() *ConversationSrv {
	conv := &ConversationSrv{
		cache:      global.Redis,
		batchStore: ext.NewBatchProcessor(1, 100, 3*time.Second, batchStoreConv),
	}
	// 启动批量异步存储处理，因消息顺序性要求workers=1 todo 优雅关闭
	conv.batchStore.Start()
	return conv
}

func (c *ConversationSrv) IndexUpdate(ctx context.Context, msg *types.MessageDTO, receivers []string) error {
	// 缓存会话最新消息
	c.cacheStoreLastMsg(ctx, msg)

	// 批量缓存用户会话
	err := c.cacheBatchStore(ctx, receivers, msg)
	if err != nil {
		logger.Err(ctx, err, "conv batch store cache")
		return err
	}
	// 异步存储用户会话数据
	for _, uid := range receivers {
		_ = c.batchStore.Submit(&convStoreData{UID: uid, LastMsg: msg})
	}
	return nil
}
func (c *ConversationSrv) cacheStoreLastMsg(ctx context.Context, msg *types.MessageDTO) {
	msgJson, _ := jsoniter.MarshalToString(msg)
	c.cache.Set(ctx, fmt.Sprintf(CacheConvLastMsg, msg.ConversationID), msgJson, 30*24*time.Hour)
}

func (c *ConversationSrv) cacheQueryLastMsg(ctx context.Context, convID string) (*types.MessageDTO, error) {
	msg, err := c.cache.Get(ctx, fmt.Sprintf(CacheConvLastMsg, convID)).Result()
	if err != nil {
		return nil, err
	}
	var mt types.MessageDTO
	_ = jsoniter.UnmarshalFromString(msg, &mt)
	return &mt, nil
}

func (c *ConversationSrv) cacheBatchStore(ctx context.Context, uidSlice []string, msg *types.MessageDTO) error {
	for _, partition := range util.Partition(uidSlice, 500) {
		pipe := c.cache.Pipeline()
		for _, uid := range partition {
			cacheInfoKey := c.getConvCacheKey(uid, msg.ConversationID)
			// 未读消息数
			if uid != msg.FromUID {
				pipe.HIncrBy(ctx, cacheInfoKey, ConvFieldUnreadCount, 1)
			}
			// 会话最新消息
			fields := map[string]interface{}{
				ConvFieldLastMsgSeq: msg.MsgSeq,
			}
			if msg.MsgSeq == 1 {
				fields[ConvFieldCreatedAt] = msg.CreatedAt
			}
			pipe.HMSet(ctx, cacheInfoKey, fields)

			// 用户最新会话
			pipe.ZAdd(ctx, fmt.Sprintf(CacheConvList, uid), redis.Z{Score: float64(msg.CreatedAt), Member: msg.ConversationID})
		}
		_, err := pipe.Exec(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *ConversationSrv) cacheQueryRange(ctx context.Context, uid string, startScore, limit int64) ([]*types.ConvDetailDTO, error) {
	opt := &redis.ZRangeBy{
		Min:   "0",
		Max:   "+inf",
		Count: limit,
	}
	if startScore > 0 {
		opt.Min = cast.ToString(startScore)
	}
	convIds, err := c.cache.ZRangeByScore(ctx, fmt.Sprintf(CacheConvList, uid), opt).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation: %v", err)
	}
	logger.Infof(ctx, "get conv ids=%v", convIds)
	if len(convIds) == 0 {
		return nil, redis.Nil
	}
	return c.cacheQueryDetail(ctx, uid, convIds)
}

func (c *ConversationSrv) cacheQueryDetail(ctx context.Context, uid string, convIds []string) ([]*types.ConvDetailDTO, error) {
	pipe := c.cache.Pipeline()
	type cmdPair struct {
		msgCmd  *redis.StringCmd
		infoCmd *redis.SliceCmd
	}
	cmds := make([]cmdPair, len(convIds))
	for i, cid := range convIds {
		cmds[i].msgCmd = pipe.Get(ctx, fmt.Sprintf(CacheConvLastMsg, cid))
		cmds[i].infoCmd = pipe.HMGet(ctx, c.getConvCacheKey(uid, cid), ConvFieldUnreadCount, ConvFieldPin, ConvFieldMute,
			ConvFieldDeleted, ConvFieldCreatedAt)
	}
	_, err := pipe.Exec(ctx)
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, fmt.Errorf("failed to execute pipeline: %v", err)
	}

	var list []*types.ConvDetailDTO
	for i := range convIds {
		conv := &types.ConvDetailDTO{UID: uid}

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
		if len(fields) == 5 {
			if fields[0] != nil {
				conv.UnreadCount = fields[0].(uint64)
			}
			if fields[1] != nil {
				conv.Pin = fields[1].(int8)
			}
			if fields[2] != nil {
				conv.Mute = fields[2].(int8)
			}
			if fields[3] != nil {
				conv.Deleted = fields[3].(int8)
			}
			if fields[4] != nil {
				conv.CreatedAt = fields[4].(int64)
			}
		}
		list = append(list, conv)
	}
	return list, nil
}

func (c *ConversationSrv) getConvCacheKey(uid, convID string) string {
	return fmt.Sprintf(CacheConvInfo, uid, convID)
}

func (c *ConversationSrv) IncrQuery(ctx context.Context, uid string, startTime, limit int64) ([]*types.ConvDetailDTO, error) {
	if list, err := c.cacheQueryRange(ctx, uid, startTime, limit); err == nil {
		return list, nil
	} else {
		logger.Err(ctx, err, fmt.Sprintf("conversation recent cache query: uid=%s", uid))
	}

	// todo 可以全用缓存，不查db
	data, err := new(store.Conversation).ListByTime(ctx, uid, startTime, limit)
	if err != nil {
		return nil, err
	}

	list := make([]*types.ConvDetailDTO, 0, len(data))
	msgIds := make([]string, 0, len(data))
	for _, d := range data {
		list = append(list, &types.ConvDetailDTO{
			UID:              uid,
			ConversationID:   d.ConversationID,
			ConversationType: d.ConversationType,
			UnreadCount:      d.UnreadCount,
			Mute:             d.Mute,
			Pin:              d.Pin,
			Deleted:          d.Deleted,
			LastMsgID:        d.LastMsgId,
			CreatedAt:        d.CreatedAt.UnixMilli(),
			Version:          d.UpdatedAt.UnixMilli(),
		})
		msgIds = append(msgIds, d.LastMsgId)
	}

	if len(msgIds) == 0 {
		return list, nil
	}
	msgList, err := new(store.Message).ListByMsgIds(ctx, msgIds)
	if err != nil {
		return nil, err
	}
	msgMap := make(map[string]*types.MessageDTO, len(msgList))
	for _, m := range msgList {
		var md types.MessageDTO
		_ = copier.Copy(&md, m)
		md.SendTime = m.SendTime.UnixMilli()
		md.CreatedAt = m.CreatedAt.UnixMilli()
		md.UpdatedAt = m.UpdatedAt.UnixMilli()
		msgMap[m.MsgID] = &md
	}
	for _, cv := range list {
		if md, ok := msgMap[cv.LastMsgID]; ok {
			cv.LastMsg = md
			cv.LastMsgSeq = md.MsgSeq
		}
	}
	return list, nil
}

func (c *ConversationSrv) SetAttribute(ctx context.Context, attr types.ConvAttributeDTO) (err error) {
	cacheKey := c.getConvCacheKey(attr.UID, attr.ConversationID)
	dbStore := new(store.Conversation)
	switch attr.Attribute {
	case ConvFieldPin:
		c.cache.HSet(ctx, cacheKey, ConvFieldPin, attr.Pin)
		err = dbStore.UpdatePin(ctx, attr.UID, attr.ConversationID, attr.Pin)
	case ConvFieldMute:
		c.cache.HSet(ctx, cacheKey, ConvFieldMute, attr.Mute)
		err = dbStore.UpdateMute(ctx, attr.UID, attr.ConversationID, attr.Mute)
	case ConvFieldUnreadCount:
		c.cache.HSet(ctx, cacheKey, ConvFieldUnreadCount, attr.UnreadCount)
		err = dbStore.UpdateUnreadCount(ctx, attr.UID, attr.ConversationID, attr.UnreadCount)
	case ConvFieldDeleted:
		c.cache.HSet(ctx, cacheKey, ConvFieldDeleted, attr.Deleted)
		err = dbStore.UpdateDeleted(ctx, attr.UID, attr.ConversationID, attr.Deleted)
	default:
		return errors.New("invalid op type")
	}
	if err != nil {
		return err
	}
	c.cache.ZAdd(ctx, fmt.Sprintf(CacheConvList, attr.UID), redis.Z{Score: float64(time.Now().UnixMilli()), Member: attr.ConversationID})
	return
}

func batchStoreConv(dataList []*convStoreData) error {
	uniqMap := make(map[string]*convStoreData)
	for _, data := range dataList {
		data.Count = 1
		if data.LastMsg.FromUID == data.UID {
			data.Count = 0
		}
		// 同一会话的连续数据，只保留最新的一条
		key := fmt.Sprintf("%s:%s", data.UID, data.LastMsg.ConversationID)
		if ed, exists := uniqMap[key]; exists {
			if data.LastMsg.MsgSeq > ed.LastMsg.MsgSeq {
				data.Count += ed.Count
				uniqMap[key] = data
			}
		} else {
			uniqMap[key] = data
		}
	}

	var bulkWrites []mongo.WriteModel
	for _, data := range uniqMap {
		// 构造查询条件
		filter := bson.M{
			"uid":             data.UID,
			"conversation_id": data.LastMsg.ConversationID,
		}
		// 构造更新操作
		update := bson.M{
			"$setOnInsert": bson.M{ // 仅在插入时设置的字段
				"conversation_id":   data.LastMsg.ConversationID,
				"conversation_type": data.LastMsg.ConversationType,
				"uid":               data.UID,
				"created_at":        time.Now(),
			},
			"$set": bson.M{ // 更新字段
				"last_msg_id":  data.LastMsg.MsgID,
				"last_msg_seq": data.LastMsg.MsgSeq,
				"from_uid":     data.LastMsg.FromUID,
				"updated_at":   time.UnixMilli(data.LastMsg.CreatedAt),
			},
			"$inc": bson.M{ // 未读计数
				"unread_count": data.Count,
			},
		}
		// 使用 Upsert 操作（如果存在则更新，否则插入）
		bulkWrite := mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update).SetUpsert(true)
		bulkWrites = append(bulkWrites, bulkWrite)
	}

	// 执行批量操作
	if len(bulkWrites) > 0 {
		_, err := new(store.Conversation).Collection().BulkWrite(context.Background(), bulkWrites, options.BulkWrite())
		if err != nil {
			return fmt.Errorf("bulk write failed: %v", err)
		}
	}
	return nil
}
