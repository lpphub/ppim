package types

type MessageDTO struct {
	ToID             string
	FromID           string
	ConversationType string
	ConversationID   string
	MsgID            string
	MsgSeq           uint64
	MsgType          int32
	Content          string
}

type OnlineDTO struct {
	Uid   string
	Did   string
	Ip    string
	Topic string
}
