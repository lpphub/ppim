package store

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"ppim/internal/logic/global"
	"time"
)

type Conversation struct {
	ConversationID   string    `bson:"conversation_id"`
	ConversationType string    `bson:"conversation_type"`
	UID              string    `bson:"uid"`
	Pin              bool      `bson:"pin"`          // 置顶
	Mute             bool      `bson:"mute"`         // 免打扰
	UnreadCount      uint64    `bson:"unread_count"` // 未读消息数
	LastMsgId        string    `bson:"last_msg_id"`
	LastMsgSeq       uint64    `bson:"last_msg_seq"`
	FromUID          string    `bson:"from_uid"`
	CreatedAt        time.Time `bson:"created_at"`
	UpdatedAt        time.Time `bson:"updated_at"`
}

func (*Conversation) Collection() *mongo.Collection {
	return global.Mongo.Collection("conversation")
}

func (c *Conversation) GetOne(ctx context.Context, uid, conversationID string) (*Conversation, error) {
	filter := bson.D{{"conversation_id", conversationID}, {"uid", uid}}
	err := c.Collection().FindOne(ctx, filter).Decode(c)
	return c, err
}

func (c *Conversation) Insert(ctx context.Context) error {
	_, err := c.Collection().InsertOne(ctx, c)
	return err
}

func (c *Conversation) Update(ctx context.Context) error {
	filter := bson.D{{"conversation_id", c.ConversationID}, {"uid", c.UID}}
	_, err := c.Collection().ReplaceOne(ctx, filter, c)
	return err
}

func (c *Conversation) ListRecent(ctx context.Context, uid string) ([]Conversation, error) {
	filter := bson.D{{"uid", uid}}
	opts := options.Find().SetSort(bson.D{{"updated_at", -1}}).SetLimit(100)

	cursor, err := c.Collection().Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	var results []Conversation
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}
