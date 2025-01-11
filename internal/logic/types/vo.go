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
	ConversationID string `json:"conversationID" form:"conversationID"`
	Pin            bool   `json:"pin" form:"pin"`
	Mute           bool   `json:"mute" form:"mute"`
	UnreadCount    uint64 `json:"unreadCount" form:"unreadCount"`
	Deleted        bool   `json:"deleted" form:"deleted"`
}

type MessageQueryVO struct {
	ConversationID string `json:"conversationID" form:"conversationID"`
	StartSeq       int64  `json:"startSeq" form:"startSeq"`
	Limit          int64  `json:"limit" form:"limit"`
}
