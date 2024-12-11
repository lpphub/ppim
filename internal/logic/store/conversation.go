package store

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"ppim/internal/logic/global"
	"time"
)

type Conversation struct {
	ConversationID   string    `bson:"conversation_id"`
	ConversationType string    `bson:"conversation_type"`
	UID              string    `bson:"uid"`
	Pin              bool      `bson:"pin"`          // 是否置顶
	NoDisturb        bool      `bson:"no_disturb"`   // 是否免打扰
	UnreadCount      uint64    `bson:"unread_count"` // 未读消息数
	LastMsgId        string    `bson:"last_msg_id"`
	LastMsgSeq       uint64    `bson:"last_msg_seq"`
	FromID           string    `bson:"from_id"`
	CreatedAt        time.Time `bson:"created_at"`
	UpdatedAt        time.Time `bson:"updated_at"`
}

func (*Conversation) Collection() *mongo.Collection {
	return global.Mongo.Collection("conversation")
}

func (c *Conversation) GetOne(ctx context.Context, uid, conversationID string) (*Conversation, error) {
	filter := bson.D{{"uid", uid}, {"conversation_id", conversationID}}
	err := c.Collection().FindOne(ctx, filter).Decode(c)
	return c, err
}

func (c *Conversation) Insert(ctx context.Context) error {
	_, err := c.Collection().InsertOne(ctx, c)
	return err
}

func (c *Conversation) Update(ctx context.Context) error {
	filter := bson.D{{"uid", c.UID}, {"conversation_id", c.ConversationID}}
	_, err := c.Collection().ReplaceOne(ctx, filter, c)
	return err
}
