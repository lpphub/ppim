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
	user UserHandler
	chat ChatHandler
)

func initHandlers() {
	user = UserHandler{srv: srv.NewUserSrv()}
	chat = ChatHandler{conv: srv.NewConvSrv(), msg: srv.NewMsgSrv()}
}

func setupRoutes(r *gin.Engine) {
	initHandlers()

	r.GET("/test", Test)

	r.Use(authorize())

	u := r.Group("/user")
	{
		u.GET("/get", user.GetOne)
		u.POST("/register", user.Register)
	}

	c := r.Group("/conversation")
	{
		c.GET("/recent", chat.RecentConvList) // 最近会话列表
		c.GET("/message", chat.ListConvMsg)   // 会话历史消息
	}
}

func Test(ctx *gin.Context) {
	t := global.Redis.Get(ctx, "test").String()
	logx.Infof(ctx, "redis test: %s", t)
	logx.Err(ctx, errors.New("bbb"), "")
	web.JsonWithSuccess(ctx, "")
}
