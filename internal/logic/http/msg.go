package http

import (
	"github.com/gin-gonic/gin"
	"github.com/lpphub/golib/render"
	"github.com/lpphub/golib/zlog"
	"ppim/internal/logic/global"
)

type MsgHandler struct {
}

func (h *MsgHandler) Test(ctx *gin.Context) {
	t := global.Redis.Get(ctx, "test").String()
	zlog.Infof(ctx, "redis test: %s", t)

	render.JsonWithSuccess(ctx, "")
}
