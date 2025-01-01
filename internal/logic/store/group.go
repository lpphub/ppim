package store

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"ppim/internal/logic/global"
	"time"
)

type Group struct {
	GID       string    `bson:"group_id"`
	UID       string    `bson:"uid"`
	CreatedAt time.Time `bson:"created_at"`
}

func (*Group) Collection() *mongo.Collection {
	return global.Mongo.Collection("group")
}

func (g *Group) ListMembers(ctx context.Context, groupId string) ([]string, error) {
	filter := bson.D{bson.E{Key: "group_id", Value: groupId}}
	opts := options.Find().SetProjection(bson.D{bson.E{Key: "uid", Value: 1}})

	cur, err := g.Collection().Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	var result []string
	if err = cur.All(context.TODO(), &result); err != nil {
		return nil, err
	}
	return result, nil
}
