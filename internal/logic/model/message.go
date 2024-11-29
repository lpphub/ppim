package model

import (
	"go.mongodb.org/mongo-driver/mongo"
	"ppim/internal/logic/global"
	"time"
)

type Message struct {
	MsgId     string    `bson:"msg_id"`
	MsgType   int       `bson:"msg_type"`
	Content   string    `bson:"content"`
	RoomId    string    `bson:"room_id"`
	FromUid   string    `bson:"from_uid"`
	CreatedAt time.Time `bson:"created_at"`
}

func (*Message) Collection() *mongo.Collection {
	return global.Mongo.Collection("message")
}
