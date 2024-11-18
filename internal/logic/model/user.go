package model

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"ppim/internal/logic/global"
)

type User struct {
	UID    string `bson:"uid"`
	DID    string `bson:"did"`
	Token  string `bson:"token"`
	Name   string `bson:"name"`
	Avatar string `bson:"avatar"`
}

func (u *User) useCollection() *mongo.Collection {
	return global.Mongo.Collection("user")
}

func (u *User) GetOne(ctx context.Context, uid string) error {
	filter := bson.D{{"uid", uid}}
	err := u.useCollection().FindOne(ctx, filter).Decode(u)
	return err
}

func (u *User) Insert(ctx context.Context) error {
	_, err := u.useCollection().InsertOne(ctx, u)
	return err
}
