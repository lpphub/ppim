package types

type UserVO struct {
	UID    string `json:"uid"`
	DID    string `json:"did"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Token  string `json:"token,omitempty"`
}

type ConvOpVO struct {
	UID            string `json:"uid" form:"uid"`
	ConversationID string `json:"conversationID" form:"conversationID"`
	Pin            bool   `json:"pin" form:"pin"`
	Mute           bool   `json:"mute" form:"mute"`
	UnreadCount    uint64 `json:"unreadCount" form:"unreadCount"`
}

type ConvMessageVO struct {
	ConversationID string `json:"conversationID" form:"conversationID"`
	StartSeq       int64  `json:"startSeq" form:"startSeq"`
	Limit          int64  `json:"limit" form:"limit"`
}
