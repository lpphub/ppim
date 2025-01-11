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
		u.GET("/get", user.Get)
		u.POST("/register", user.Register)
	}

	c := r.Group("/conversation")
	{
		c.GET("/list", chat.ConvList)             // 会话列表
		c.GET("/message", chat.ConvListMsg)       // 会话历史消息
		c.POST("/pin", chat.ConvPin)              // 会话置顶
		c.POST("/mute", chat.ConvMute)            // 会话免打扰
		c.POST("/set_unread", chat.ConvSetUnread) // 会话未读数
		c.DELETE("/del", chat.ConvDel)            // 会话删除
	}

	m := r.Group("/message")
	{
		m.POST("/revoke", chat.MsgRevoke) // 撤回消息
		m.DELETE("/del", chat.MsgDel)     // 删除消息
		m.POST("/read", chat.MsgRead)     // 消息已读回执
	}
}

func Test(ctx *gin.Context) {
	t := global.Redis.Get(ctx, "test").String()
	logx.Infof(ctx, "redis test: %s", t)
	logx.Err(ctx, errors.New("bbb"), "")
	web.JsonWithSuccess(ctx, "")
}
