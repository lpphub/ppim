package chatlib

import (
	jsoniter "github.com/json-iterator/go"
)

const (
	DeliveryEvent  = "event"
	DeliveryChat   = "chat"
	DeliveryNotify = "notify"
)

type DeliveryMsg struct {
	CMD       string   `json:"CMD"` // 指令：chat-聊天 event-事件 notify-通知
	FromUID   string   `json:"fromUID"`
	Receivers []string `json:"receivers"`
	Chat      *ChatMsg `json:"chatMsg,omitempty"`
}

type ChatMsg struct {
	MsgID            string `json:"msgID"`            // 消息ID (唯一标识)
	MsgSeq           uint64 `json:"msgSeq"`           // 消息序列号 (递增)
	MsgNo            string `json:"msgNo"`            // 消息编号 (客户端编号)
	ConversationType string `json:"conversationType"` // 会话类型 (单聊、群聊)
	ConversationID   string `json:"conversationID"`   // 会话ID
	MsgType          int8   `json:"msgType"`          // 消息类型 (文本、图片、语音、视频、文件等)
	Content          string `json:"content"`          // 消息内容
	ToID             string `json:"toID"`             // 接收者ID
	FromUID          string `json:"fromUID"`          // 发送者UID
	FromDID          string `json:"fromDID"`          // 发送者DID
	SendTime         int64  `json:"sendTime"`         // 发送时间
	CreatedAt        int64  `json:"createdAt"`        // 创建时间
}

func (d *DeliveryMsg) ToJsonBytes() []byte {
	bytes, _ := jsoniter.Marshal(d)
	return bytes
}
