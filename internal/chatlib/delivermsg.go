package chatlib

import (
	jsoniter "github.com/json-iterator/go"
	"ppim/internal/logic/types"
)

const (
	DeliverEvent  = "event"
	DeliverChat   = "chat"
	DeliverNotify = "notify"
)

type DeliverMsg struct {
	CMD     string // 指令：event-事件 chat-聊天 notify-通知
	ToUID   string
	ChatMsg *types.MessageDTO
}

func (d *DeliverMsg) ToJsonBytes() []byte {
	bytes, _ := jsoniter.Marshal(d)
	return bytes
}
