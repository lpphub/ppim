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
	Pin              int8      `bson:"pin"`          // 置顶
	Mute             int8      `bson:"mute"`         // 免打扰
	UnreadCount      uint64    `bson:"unread_count"` // 未读消息数
	LastMsgId        string    `bson:"last_msg_id"`
	LastMsgSeq       uint64    `bson:"last_msg_seq"`
	FromUID          string    `bson:"from_uid"`
	Deleted          int8      `bson:"deleted"`
	CreatedAt        time.Time `bson:"created_at"`
	UpdatedAt        time.Time `bson:"updated_at"`
}

func (*Conversation) Collection() *mongo.Collection {
	return global.Mongo.Collection("conversation")
}

func (c *Conversation) GetOne(ctx context.Context, uid, conversationID string) (*Conversation, error) {
	filter := bson.D{bson.E{Key: "conversation_id", Value: conversationID}, bson.E{Key: "uid", Value: uid}}
	err := c.Collection().FindOne(ctx, filter).Decode(c)
	return c, err
}

func (c *Conversation) Insert(ctx context.Context) error {
	_, err := c.Collection().InsertOne(ctx, c)
	return err
}

func (c *Conversation) Update(ctx context.Context) error {
	filter := bson.D{{Key: "conversation_id", Value: c.ConversationID}, {Key: "uid", Value: c.UID}}
	_, err := c.Collection().ReplaceOne(ctx, filter, c)
	return err
}

func (c *Conversation) ListByTime(ctx context.Context, uid string, startTime, limit int64) ([]Conversation, error) {
	filter := bson.D{{Key: "uid", Value: uid}}
	if startTime > 0 {
		filter = append(filter, bson.E{Key: "updated_at", Value: bson.M{"$gte": time.UnixMilli(startTime)}})
	}
	opts := options.Find().SetSort(bson.D{{Key: "updated_at", Value: 1}}).SetLimit(limit)

	cursor, err := c.Collection().Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	var result []Conversation
	if err = cursor.All(ctx, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Conversation) GetMaxSeq(ctx context.Context, conversationID string) (uint64, error) {
	filter := bson.D{{Key: "conversation_id", Value: conversationID}}
	opts := options.FindOne().SetProjection(bson.M{"last_msg_seq": 1}).SetSort(bson.D{{Key: "last_msg_seq", Value: -1}})

	err := c.Collection().FindOne(ctx, filter, opts).Decode(c)
	if err != nil {
		return 0, err
	}
	return c.LastMsgSeq, nil
}

func (c *Conversation) UpdatePin(ctx context.Context, uid, conversationID string, pin int8) error {
	return c.UpdateField(ctx, uid, conversationID, "pin", pin)
}

func (c *Conversation) UpdateMute(ctx context.Context, uid, conversationID string, mute int8) error {
	return c.UpdateField(ctx, uid, conversationID, "mute", mute)
}

func (c *Conversation) UpdateUnreadCount(ctx context.Context, uid, conversationID string, unreadCount uint64) error {
	return c.UpdateField(ctx, uid, conversationID, "unread_count", unreadCount)
}

func (c *Conversation) UpdateDeleted(ctx context.Context, uid, conversationID string, deleted int8) error {
	return c.UpdateField(ctx, uid, conversationID, "deleted", deleted)
}

func (c *Conversation) UpdateField(ctx context.Context, uid, conversationID string, field string, value interface{}) error {
	filter := bson.D{{Key: "conversation_id", Value: conversationID}, {Key: "uid", Value: uid}}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: field, Value: value},
		{Key: "updated_at", Value: time.Now()},
	}}}
	_, err := c.Collection().UpdateOne(ctx, filter, update)
	return err
}
