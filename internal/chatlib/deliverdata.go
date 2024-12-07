package chatlib

type DeliverData struct {
	CMD     string // 指令：event-事件 msg-消息 notify-通知
	ToUID   []string
	MsgData []byte
}
