package model

import (
	"go.mongodb.org/mongo-driver/mongo"
	"ppim/internal/logic/global"
	"time"
)

type Room struct {
	RooId     string    `bson:"room_id"`
	RoomType  int       `bson:"room_type"`
	CreatedAt time.Time `bson:"created_at"`
}

type RoomUser struct {
	RoomId      string    `bson:"room_id"`
	Uid         string    `bson:"uid"`
	IsTop       bool      `bson:"is_top"`        // 置顶
	IsNoDisturb bool      `bson:"is_no_disturb"` // 免打扰
	CreatedAt   time.Time `bson:"created_at"`
}

func (*Room) Collection() *mongo.Collection {
	return global.Mongo.Collection("room")
}

func (*RoomUser) Collection() *mongo.Collection {
	return global.Mongo.Collection("room_user")
}
