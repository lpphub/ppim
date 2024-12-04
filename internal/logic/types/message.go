package types

import "time"

type MessageDTO struct {
	MsgID            string    // 消息ID (唯一标识)
	MsgSeq           uint64    // 消息序列号 (递增)
	MsgNo            string    // 消息编号 (客户端编号)
	ConversationType string    // 会话类型 (单聊、群聊)
	ConversationID   string    // 会话ID
	MsgType          int32     // 消息类型 (文本、图片、语音、视频、文件等)
	Content          string    // 消息内容
	ToID             string    // 接收者ID
	FromID           string    // 发送者ID
	CreatedAt        time.Time // 创建时间
}

type OnlineDTO struct {
	Uid   string
	Did   string
	Ip    string
	Topic string
}
