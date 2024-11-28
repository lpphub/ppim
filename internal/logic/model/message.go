package model

import (
	"go.mongodb.org/mongo-driver/mongo"
	"ppim/internal/logic/global"
	"time"
)

type Message struct {
	msgId     string    `bson:"msg_id"`
	msgType   int       `bson:"msg_type"`
	content   string    `bson:"content"`
	seqNo     string    `bson:"seq_no"`
	fromUid   string    `bson:"from_uid"`
	createdAt time.Time `bson:"created_at"`
}

type MessageRecord struct {
	MsgId     string    `bson:"msg_id"`
	RoomId    string    `bson:"room_id"`
	CreatedAt time.Time `bson:"created_at"`
}

func (*Message) Collection() *mongo.Collection {
	return global.Mongo.Collection("message")
}

func (*MessageRecord) Collection() *mongo.Collection {
	return global.Mongo.Collection("message_record")
}
