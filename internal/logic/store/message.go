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
	MsgID            string    `bson:"msg_id"`
	MsgSeq           uint64    `bson:"msg_seq"`
	MsgNo            string    `bson:"msg_no"`
	MsgType          int8      `bson:"msg_type"`
	Content          string    `bson:"content"`
	ConversationID   string    `bson:"conversation_id"`
	ConversationType string    `bson:"conversation_type"`
	FromUID          string    `bson:"from_uid"`
	ToID             string    `bson:"to_id"`
	SendTime         time.Time `bson:"send_time"`
	CreatedAt        time.Time `bson:"created_at"`
	UpdatedAt        time.Time `bson:"updated_at"`
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

func (m *Message) ListByConvSeq(ctx context.Context, conversationID string, startSeq, limit int64) ([]Message, error) {
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
