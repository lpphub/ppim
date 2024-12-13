package rpctypes

type AuthReq struct {
	Uid   string
	Did   string
	Token string
}

type AuthResp struct {
	Code int
	Msg  string
}

type RouterReq struct {
	Uid   string
	Did   string
	Topic string
	Ip    string
}

type RouterResp struct {
}

type MessageReq struct {
	ToID             string
	FromID           string
	ConversationType string
	ConversationID   string
	MsgID            string
	MsgSeq           uint64
	MsgNo            string
	MsgType          int32
	Content          string
	SendTime         uint64
}

type MessageResp struct {
}
