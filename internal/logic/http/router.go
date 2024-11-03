package http

import (
	"github.com/gin-gonic/gin"
	"ppim/internal/logic/srv"
)

var (
	user *UserHandler
	msg  *MsgHandler
)

func initSvcCtx() {
	user = &UserHandler{Srv: srv.NewUserSrv()}
	msg = &MsgHandler{}
}

func registerRoute(r *gin.Engine) {
	initSvcCtx()

	r.GET("/test", msg.Test)

	u := r.Group("/user")
	{
		u.GET("/get", user.GetOne)
		u.POST("/post", user.Create)
	}
}
