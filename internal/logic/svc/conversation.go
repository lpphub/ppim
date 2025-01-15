package svc

import (
	"context"
	"errors"
	"fmt"
	"github.com/jinzhu/copier"
	jsoniter "github.com/json-iterator/go"
	"github.com/lpphub/golib/gowork"
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
	workers     *gowork.Pool
	segmentLock *ext.SegmentRWLock
	batchStore  *ext.BatchProcessor[*convStoreData]
}

type convStoreData struct {
	UID     string
	LastMsg *types.MessageDTO
}

func newConversationSrv() *ConversationSrv {
	conv := &ConversationSrv{
		workers:     gowork.NewPool(100),
		segmentLock: ext.NewSegmentLock(20),
		batchStore:  ext.NewBatchProcessor(context.Background(), 1, 100, 5*time.Second, batchStoreConv),
	}
	// 启动批量异步存储处理器，因消息顺序性只能workerCount=1 todo 优雅关闭
	conv.batchStore.Start()
	return conv
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

	recentConvMaxSize = 1000
)

func (c *ConversationSrv) IndexRecent(ctx context.Context, msg *types.MessageDTO, receivers []string) error {
	msgJson, _ := jsoniter.MarshalToString(msg)
	global.Redis.Set(ctx, fmt.Sprintf(CacheConvLastMsg, msg.ConversationID), msgJson, 30*24*time.Hour)

	for _, uid := range receivers {
		_ = c.workers.Submit(func() {
			c.indexWithLock(ctx, uid, msg)
		})
	}
	return nil
}

func (c *ConversationSrv) indexWithLock(ctx context.Context, uid string, msg *types.MessageDTO) {
	// todo 集群模式下，分布式锁
	index := cast.ToInt(util.Murmur32(fmt.Sprintf("%s-%s", uid, msg.ConversationID)))
	c.segmentLock.Lock(index)
	defer c.segmentLock.Unlock(index)

	if err := c.cacheStoreRecent(ctx, uid, msg); err != nil {
		logger.Err(ctx, err, fmt.Sprintf("conversation cache recent: uid=%s", uid))
		return
	}

	_ = c.batchStore.Submit(&convStoreData{UID: uid, LastMsg: msg})

	// 优化：异步周期性从redis持久化至存储层，减少数据存储操作
	//conv, err := new(store.Conversation).GetOne(ctx, uid, msg.ConversationID)
	//if err != nil && errors.Is(err, mongo.ErrNoDocuments) {
	//	conv = &store.Conversation{
	//		ConversationID:   msg.ConversationID,
	//		ConversationType: msg.ConversationType,
	//		UID:              uid,
	//		UnreadCount:      1,
	//		LastMsgId:        msg.MsgID,
	//		LastMsgSeq:       msg.MsgSeq,
	//		CreatedAt:        time.Now(),
	//		UpdatedAt:        time.Now(),
	//	}
	//	_ = conv.Insert(ctx)
	//} else {
	//	if msg.FromUID != uid {
	//		conv.UnreadCount++
	//	}
	//	if conv.LastMsgSeq < msg.MsgSeq {
	//		conv.LastMsgId = msg.MsgID
	//		conv.LastMsgSeq = msg.MsgSeq
	//		conv.FromUID = msg.FromUID
	//	}
	//	conv.UpdatedAt = time.UnixMilli(msg.CreatedAt)
	//	_ = conv.Update(ctx)
	//}
}

func (c *ConversationSrv) getConvCacheKey(uid, convID string) string {
	return fmt.Sprintf(CacheConvInfo, uid, convID)
}

// 缓存用户最近会话
func (c *ConversationSrv) cacheStoreRecent(ctx context.Context, uid string, msg *types.MessageDTO) error {
	cacheKey := fmt.Sprintf(CacheConvList, uid)
	cacheInfoKey := c.getConvCacheKey(uid, msg.ConversationID)

	rdb := global.Redis
	pipe := rdb.Pipeline()
	// 会话最新消息
	pipe.HSet(ctx, cacheInfoKey, ConvFieldLastMsgSeq, msg.MsgSeq)
	if uid != msg.FromUID {
		// 未读消息数
		pipe.HIncrBy(ctx, cacheInfoKey, ConvFieldUnreadCount, 1)
	}
	if msg.MsgSeq == 1 {
		pipe.HSet(ctx, cacheInfoKey, ConvFieldCreatedAt, msg.CreatedAt)
	}
	// 用户最新会话
	pipe.ZAdd(ctx, cacheKey, redis.Z{Score: float64(msg.CreatedAt), Member: msg.ConversationID})
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

func (c *ConversationSrv) cacheQueryRange(ctx context.Context, uid string, startScore, limit int64) ([]*types.ConvDetailDTO, error) {
	opt := &redis.ZRangeBy{
		Min:   "0",
		Max:   "+inf",
		Count: limit,
	}
	if startScore > 0 {
		opt.Min = cast.ToString(startScore)
	}
	convIds, err := global.Redis.ZRangeByScore(ctx, fmt.Sprintf(CacheConvList, uid), opt).Result()
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
	pipe := global.Redis.Pipeline()
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

func (c *ConversationSrv) ListByUID(ctx context.Context, uid string, startTime, limit int64) ([]*types.ConvDetailDTO, error) {
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
	case ConvFieldDeleted:
		global.Redis.HSet(ctx, cacheKey, ConvFieldDeleted, attr.Deleted)
		err = new(store.Conversation).UpdateDeleted(ctx, attr.UID, attr.ConversationID, attr.Deleted)
	default:
		return errors.New("invalid op type")
	}
	if err != nil {
		return err
	}
	global.Redis.ZAdd(ctx, fmt.Sprintf(CacheConvList, attr.UID), redis.Z{Score: float64(time.Now().UnixMilli()), Member: attr.ConversationID})
	return
}

func batchStoreConv(dataList []*convStoreData) error {
	uniqMap := make(map[string]*convStoreData)
	for _, data := range dataList {
		key := fmt.Sprintf("%s:%s", data.UID, data.LastMsg.ConversationID)
		if ed, exists := uniqMap[key]; exists {
			if data.LastMsg.MsgSeq > ed.LastMsg.MsgSeq {
				uniqMap[key] = data
			}
		} else {
			uniqMap[key] = data
		}
	}

	var bulkWrites []mongo.WriteModel
	ctx := context.Background()

	for _, data := range uniqMap {
		unReadCount := 1
		if data.LastMsg.FromUID == data.UID {
			unReadCount = 0
		}

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
				"unread_count": unReadCount,
			},
		}
		// 使用 Upsert 操作（如果存在则更新，否则插入）
		bulkWrite := mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update).SetUpsert(true)
		bulkWrites = append(bulkWrites, bulkWrite)
	}

	// 执行批量操作
	if len(bulkWrites) > 0 {
		_, err := new(store.Conversation).Collection().BulkWrite(ctx, bulkWrites, options.BulkWrite())
		if err != nil {
			return fmt.Errorf("bulk write failed: %v", err)
		}
	}
	return nil
}
