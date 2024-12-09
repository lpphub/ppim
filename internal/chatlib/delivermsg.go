package chatlib

import jsoniter "github.com/json-iterator/go"

type DeliverMsg struct {
	CMD     string // 指令：event-事件 chat-聊天 notify-通知
	ToUID   []string
	Content []byte
}

func (d *DeliverMsg) ToJsonBytes() []byte {
	bytes, _ := jsoniter.Marshal(d)
	return bytes
}
