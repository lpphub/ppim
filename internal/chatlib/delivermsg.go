package chatlib

import (
	jsoniter "github.com/json-iterator/go"
	"ppim/internal/logic/types"
)

const (
	DeliveryEvent  = "event"
	DeliveryChat   = "chat"
	DeliveryNotify = "notify"
)

type DeliveryMsg struct {
	CMD       string // 指令：event-事件 chat-聊天 notify-通知
	Receivers []string
	ChatMsg   *types.MessageDTO
}

func (d *DeliveryMsg) ToJsonBytes() []byte {
	bytes, _ := jsoniter.Marshal(d)
	return bytes
}
