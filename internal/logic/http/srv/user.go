package srv

import (
	"github.com/gin-gonic/gin"
	"github.com/lpphub/golib/logger/logx"
	"ppim/internal/logic/store"
	"ppim/internal/logic/types"
)

type UserSrv struct {
}

func NewUserSrv() *UserSrv {
	return &UserSrv{}
}

func (srv *UserSrv) GetOne(ctx *gin.Context, uid string) (resp *types.UserDTO, err error) {
	logx.Infof(ctx, "uid: %s", uid)
	user, err := new(store.User).GetOne(ctx, uid)
	if err != nil {
		return
	}
	resp = &types.UserDTO{
		UID:    uid,
		DID:    user.DID,
		Name:   user.Name,
		Avatar: user.Avatar,
	}
	return
}

func (srv *UserSrv) Register(ctx *gin.Context, req types.UserDTO) error {
	logx.Infof(ctx, "user register param: %v", req)
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
