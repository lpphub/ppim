package store

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

func (u *User) Collection() *mongo.Collection {
	return global.Mongo.Collection("user")
}

func (u *User) GetOne(ctx context.Context, uid string) (*User, error) {
	filter := bson.D{{"uid", uid}}
	err := u.Collection().FindOne(ctx, filter).Decode(u)
	return u, err
}

func (u *User) Insert(ctx context.Context) error {
	_, err := u.Collection().InsertOne(ctx, u)
	return err
}
