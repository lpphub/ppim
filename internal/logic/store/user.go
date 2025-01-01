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

func (u *User) List(ctx context.Context, uid string) ([]User, error) {
	filter := bson.D{{Key: "uid", Value: uid}}
	cursor, err := u.Collection().Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var result []User
	if err = cursor.All(ctx, &result); err != nil {
		return nil, err
	}
	return result, err
}

func (u *User) GetByToken(ctx context.Context, uid, token string) (*User, error) {
	filter := bson.D{{Key: "uid", Value: uid}, {Key: "token", Value: token}}
	err := u.Collection().FindOne(ctx, filter).Decode(u)
	return u, err
}

func (u *User) Insert(ctx context.Context) error {
	_, err := u.Collection().InsertOne(ctx, u)
	return err
}
