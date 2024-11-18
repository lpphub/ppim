package types

type UserDTO struct {
	UID    string `json:"uid"`
	DID    string `json:"did"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Token  string `json:"token,omitempty"`
}
