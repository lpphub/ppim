package model

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"ppim/internal/logic/global"
	"time"
)

type Message struct {
	MsgId            string    `bson:"msg_id"`
	MsgType          int32     `bson:"msg_type"`
	Content          string    `bson:"content"`
	ConversationType string    `bson:"conversation_type"`
	FromID           string    `bson:"from_id"`
	ToID             string    `bson:"to_id"`
	CreatedAt        time.Time `bson:"created_at"`
}

func (*Message) Collection() *mongo.Collection {
	return global.Mongo.Collection("message")
}

func (m *Message) Insert(ctx context.Context) error {
	_, err := m.Collection().InsertOne(ctx, m)
	return err
}
