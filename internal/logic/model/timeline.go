package model

import (
	"go.mongodb.org/mongo-driver/mongo"
	"ppim/internal/logic/global"
	"time"
)

type Timeline struct {
	Uid       string    `bson:"uid"`
	RoomId    string    `bson:"room_id"`
	MsgId     string    `bson:"msg_id"`
	IsRead    bool      `bson:"is_read"`
	CreatedAt time.Time `bson:"created_at"`
}

func (*Timeline) Collection() *mongo.Collection {
	return global.Mongo.Collection("timeline")
}
