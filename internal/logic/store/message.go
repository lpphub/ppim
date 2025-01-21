package store

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"ppim/internal/logic/global"
	"time"
)

type Message struct {
	MsgID            string     `bson:"msg_id"`            //消息id
	MsgSeq           uint64     `bson:"msg_seq"`           //消息序号(server)
	MsgNo            string     `bson:"msg_no"`            //消息编号(client)
	MsgType          int8       `bson:"msg_type"`          //消息类型
	Content          string     `bson:"content"`           //消息内容
	ConversationID   string     `bson:"conversation_id"`   //会话id
	ConversationType string     `bson:"conversation_type"` //会话类型
	FromUID          string     `bson:"from_uid"`          //发送者
	ToID             string     `bson:"to_id"`             //接收者
	Revoked          int8       `bson:"revoked"`           //撤回
	Deleted          int8       `bson:"deleted"`           //删除
	RevokedTime      *time.Time `bson:"revoked_time"`      //撤回时间
	DeletedTime      *time.Time `bson:"deleted_time"`      //删除时间
	SendTime         time.Time  `bson:"send_time"`         //发送时间
	CreatedAt        time.Time  `bson:"created_at"`
	UpdatedAt        time.Time  `bson:"updated_at"`
}

func (*Message) Collection() *mongo.Collection {
	return global.Mongo.Collection("message")
}

func (m *Message) Insert(ctx context.Context) error {
	_, err := m.Collection().InsertOne(ctx, m)
	return err
}

func (m *Message) ListByMsgIds(ctx context.Context, msgIds []string) ([]Message, error) {
	filter := bson.M{"msg_id": bson.M{"$in": msgIds}}
	cursor, err := m.Collection().Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var result []Message
	if err = cursor.All(ctx, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// PullUpOrDown limit: 负数时向前-10: 90~99, 正数时向后10: 101~110
func (m *Message) PullUpOrDown(ctx context.Context, conversationID string, startSeq, limit int64) ([]Message, error) {
	filter := bson.D{{Key: "conversation_id", Value: conversationID}}
	opts := options.Find()

	if limit < 0 { // 向前翻页
		filter = append(filter, bson.E{Key: "msg_seq", Value: bson.M{"$lt": startSeq}})
		opts.SetSort(bson.D{{Key: "msg_seq", Value: -1}}).SetLimit(limit * -1)
	} else { // 向后翻页
		filter = append(filter, bson.E{Key: "msg_seq", Value: bson.M{"$gt": startSeq}})
		opts.SetSort(bson.D{{Key: "msg_seq", Value: 1}}).SetLimit(limit)
	}

	cursor, err := m.Collection().Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	var result []Message
	if err = cursor.All(ctx, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (m *Message) Revoke(ctx context.Context, msgID string) error {
	fields := map[string]interface{}{
		"revoked":      1,
		"revoked_time": time.Now(),
	}
	return m.UpdateFields(ctx, msgID, fields)
}

func (m *Message) Delete(ctx context.Context, msgID string) error {
	fields := map[string]interface{}{
		"deleted":      1,
		"deleted_time": time.Now(),
	}
	return m.UpdateFields(ctx, msgID, fields)
}

func (m *Message) UpdateFields(ctx context.Context, msgID string, fieldMap map[string]interface{}) error {
	fields := bson.D{bson.E{Key: "updated_at", Value: time.Now()}}
	for k, v := range fieldMap {
		fields = append(fields, bson.E{Key: k, Value: v})
	}
	filter := bson.M{"msg_id": msgID}
	update := bson.D{{Key: "$set", Value: fields}}
	_, err := m.Collection().UpdateOne(ctx, filter, update)
	return err
}
