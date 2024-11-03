package router

import (
	"github.com/gin-gonic/gin"
	"ppim/internal/logic/router/handler"
	"ppim/internal/logic/srv"
)

var (
	user *handler.UserHandler
)

func init() {
	user = &handler.UserHandler{Srv: srv.NewUserSrv()}
}

func registerRoute(r *gin.Engine) {
	r.GET("/test", user.Test)

	u := r.Group("/user")
	{
		u.GET("/get", user.GetOne)
		u.POST("/post", user.Create)
	}
}
