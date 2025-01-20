package types

type UserVO struct {
	UID    string `json:"uid"`
	DID    string `json:"did"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Token  string `json:"token,omitempty"`
}

type ConvQueryVO struct {
	UID     string `json:"uid" form:"uid"`
	NextKey string `json:"nextKey" form:"nextKey"` // 增量标识
	Limit   int    `json:"limit" form:"limit"`
}

type ConvListVO struct {
	NextKey string           `json:"nextKey" form:"nextKey"`
	List    []*ConvDetailDTO `json:"list"`
}

type ConvOpVO struct {
	UID            string `json:"uid" form:"uid"`
	ConversationID string `json:"conversationID" form:"conversationID" binding:"required"`
	Pin            int8   `json:"pin" form:"pin"`
	Mute           int8   `json:"mute" form:"mute"`
	UnreadCount    uint64 `json:"unreadCount" form:"unreadCount"`
	Deleted        int8   `json:"deleted" form:"deleted"`
}

type MessageQueryVO struct {
	ConversationID string `json:"conversationID" form:"conversationID" binding:"required"`
	StartSeq       int64  `json:"startSeq" form:"startSeq"`
	Limit          int64  `json:"limit" form:"limit"`
}

type MsgOpVO struct {
	ConversationID string `json:"conversationID" form:"conversationID" binding:"required"`
	MsgID          string `json:"msgID" form:"msgID" binding:"required"`
}
