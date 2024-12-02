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
	Ip    string
	Topic string
}

type RouterResp struct {
}

type MessageReq struct {
	ToID             string
	FromID           string
	ConversationType string
	ConversationID   string
	MsgID            string
	Sequence         uint64
	MsgType          int32
	Content          string
}

type MessageResp struct {
	Sequence uint64
}
