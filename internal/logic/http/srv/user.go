package srv

import (
	"github.com/gin-gonic/gin"
	"github.com/lpphub/golib/logger/glog"
	"ppim/internal/logic/model"
	"ppim/internal/logic/types"
)

type UserSrv struct {
}

func NewUserSrv() *UserSrv {
	return &UserSrv{}
}

func (srv *UserSrv) GetOne(ctx *gin.Context, uid string) (resp *types.UserDTO, err error) {
	glog.Infof(ctx, "uid: %s", uid)
	user := &model.User{}
	err = user.GetOne(ctx, uid)
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
	glog.Infof(ctx, "user register param: %v", req)
	doc := model.User{
		UID:    req.UID,
		DID:    req.DID,
		Token:  req.Token,
		Name:   req.Name,
		Avatar: req.Avatar,
	}
	// todo 缓存token

	return doc.Insert(ctx)
}
