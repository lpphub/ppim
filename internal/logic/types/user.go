package types

type UserGetOneResp struct {
	UID    int64  `json:"uid"`
	Name   string `json:"name"`
	Mobile string `json:"mobile"`
	Email  string `json:"email"`
}

type UserCreateReq struct {
	Name   string `json:"name"`
	Mobile string `json:"mobile"`
	Email  string `json:"email"`
}
