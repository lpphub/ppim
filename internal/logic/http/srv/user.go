package srv

import (
	"github.com/gin-gonic/gin"
	"ppim/internal/logic/store"
	"ppim/internal/logic/types"
)

type UserSrv struct {
}

func NewUserSrv() *UserSrv {
	return &UserSrv{}
}

func (srv *UserSrv) Get(ctx *gin.Context, uid string) (resp []*types.UserDTO, err error) {
	users, err := new(store.User).List(ctx, uid)
	if err != nil {
		return
	}
	for _, user := range users {
		u := &types.UserDTO{
			UID:    uid,
			DID:    user.DID,
			Name:   user.Name,
			Avatar: user.Avatar,
		}
		resp = append(resp, u)
	}
	return
}

func (srv *UserSrv) Register(ctx *gin.Context, req types.UserDTO) error {
	doc := store.User{
		UID:    req.UID,
		DID:    req.DID,
		Token:  req.Token,
		Name:   req.Name,
		Avatar: req.Avatar,
	}
	// todo 缓存token

	return doc.Insert(ctx)
}
