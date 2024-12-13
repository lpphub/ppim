package http

import (
	"github.com/gin-gonic/gin"
	"github.com/lpphub/golib/logger/logx"
	"github.com/lpphub/golib/web"
	"github.com/pkg/errors"
	"ppim/internal/logic/global"
	"ppim/internal/logic/http/srv"
)

var (
	user *UserHandler
	msg  *MsgHandler
)

func initSvcCtx() {
	user = &UserHandler{Srv: srv.NewUserSrv()}
	msg = &MsgHandler{ConvSrv: srv.NewConversationSrv()}
}

func registerRoutes(r *gin.Engine) {
	initSvcCtx()

	r.GET("/test", Test)

	u := r.Group("/user")
	{
		u.GET("/get", user.GetOne)
		u.POST("/register", user.Register)
	}

	c := r.Group("/conversation")
	{
		c.GET("/recent", msg.RecentConvList)
	}
}

func Test(ctx *gin.Context) {
	t := global.Redis.Get(ctx, "test").String()
	logx.Infof(ctx, "redis test: %s", t)
	logx.Err(ctx, errors.New("bbb"), "")
	web.JsonWithSuccess(ctx, "")
}
