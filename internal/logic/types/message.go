package types

type MessageDTO struct {
	MsgID            string `json:"msgID"`            // 消息ID (唯一标识)
	MsgSeq           uint64 `json:"msgSeq"`           // 消息序列号 (递增)
	MsgNo            string `json:"msgNo"`            // 消息编号 (客户端编号)
	ConversationType string `json:"conversationType"` // 会话类型 (单聊、群聊)
	ConversationID   string `json:"conversationID"`   // 会话ID
	MsgType          int8   `json:"msgType"`          // 消息类型 (文本、图片、语音、视频、文件等)
	Content          string `json:"content"`          // 消息内容
	ToID             string `json:"toID"`             // 接收者UID
	FromUID          string `json:"fromUID"`          // 发送者UID
	FromDID          string `json:"fromDID"`          // 发送者DID
	SendTime         int64  `json:"sendTime"`         // 发送时间
	CreatedAt        int64  `json:"createdAt"`        // 创建时间
}

type RouteDTO struct {
	Uid   string
	Did   string
	Topic string
	Ip    string
}

type RecentConvVO struct {
	ConversationID   string      `json:"conversationID"`
	ConversationType string      `json:"conversationType"`
	UnreadCount      uint64      `json:"unreadCount"`
	Pin              bool        `json:"pin"`
	Mute             bool        `json:"mute"`
	FromUID          string      `json:"fromUID"`
	LastMsgID        string      `json:"lastMsgID"`
	LastMsgSeq       uint64      `json:"lastMsgSeq"`
	LastMsg          *MessageDTO `json:"lastMsg"`
	Version          int64       `json:"version"`
}

type ConvMsgDTO struct {
	ConversationID string `json:"conversationID" form:"conversationID"`
	StartSeq       int64  `json:"startSeq" form:"startSeq"`
	Limit          int64  `json:"limit" form:"limit"`
}
