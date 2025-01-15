package svc

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"ppim/internal/logic/store"
	"ppim/internal/logic/types"
	"ppim/pkg/ext"
	"time"
)

type ConvBatchData struct {
	UID     string
	LastMsg *types.MessageDTO
}

func batchStoreConv(dataList []*ConvBatchData) error {
	uniqMap := make(map[string]*ConvBatchData)
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

type ConvBatchStore struct {
	processor *ext.BatchProcessor[*ConvBatchData]
}

func newConvBatchStore(batchSize int, interval time.Duration) *ConvBatchStore {
	t := &ConvBatchStore{
		processor: ext.NewBatchProcessor(context.Background(), batchSize, interval, batchStoreConv),
	}

	t.processor.Start()
	return t
}

func (c *ConvBatchStore) Submit(data *ConvBatchData) error {
	return c.processor.Submit(data)
}
