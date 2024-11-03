package srv

import (
	"github.com/gin-gonic/gin"
	"github.com/lpphub/golib/zlog"
	"go.mongodb.org/mongo-driver/bson"
	"ppim/internal/logic/global"
	"ppim/internal/logic/model"
	"ppim/internal/logic/types"
)

type UserSrv struct {
}

func NewUserSrv() *UserSrv {
	return &UserSrv{}
}

func (srv *UserSrv) GetOne(ctx *gin.Context, uid int64) (resp *types.UserGetOneResp, err error) {
	zlog.Infof(ctx, "uid: %d", uid)

	var user model.User
	filter := bson.D{{"uid", uid}}
	err = global.Mongo.Collection("user").FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return
	}
	resp = &types.UserGetOneResp{
		UID:    uid,
		Name:   user.Name,
		Mobile: user.Mobile,
		Email:  user.Email,
	}
	return
}

func (srv *UserSrv) Create(ctx *gin.Context, req types.UserCreateReq) (uid int64, err error) {
	zlog.Infof(ctx, "user create param: %v", req)
	doc := model.User{
		UID:    global.Snowflake.Generate().Int64(),
		Name:   req.Name,
		Email:  req.Email,
		Mobile: req.Mobile,
	}
	_, err = global.Mongo.Collection("user").InsertOne(ctx, doc)
	uid = doc.UID
	return
}
