package store

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

	var results []Message
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}
