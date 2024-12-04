package store

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"ppim/internal/logic/global"
	"time"
)

type Message struct {
	MsgID            string    `bson:"msg_id"`
	MsgSeq           uint64    `bson:"msg_seq"`
	MsgNo            string    `bson:"msg_no"`
	MsgType          int32     `bson:"msg_type"`
	Content          string    `bson:"content"`
	ConversationID   string    `bson:"conversation_id"`
	ConversationType string    `bson:"conversation_type"`
	FromID           string    `bson:"from_id"`
	ToID             string    `bson:"to_id"`
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
