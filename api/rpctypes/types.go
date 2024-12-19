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
	ToID             string // 目标用户
	FromUID          string // 发送者
	FromDID          string // 发送者DID
	ConversationType string // 会话类型
	ConversationID   string // 会话ID
	MsgID            string // 消息ID
	MsgSeq           uint64 // 消息序列号
	MsgNo            string // 消息编号
	MsgType          int8   // 消息类型
	Content          string // 消息内容
	SendTime         int64  // 发送时间
	CreatedAt        int64  // 创建时间
}

type MessageResp struct {
}
