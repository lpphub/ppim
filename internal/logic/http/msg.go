package http

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/lpphub/golib/logger/logx"
	"github.com/lpphub/golib/web"
	"ppim/internal/logic/global"
)

type MsgHandler struct {
}

func (h *MsgHandler) Test(ctx *gin.Context) {
	t := global.Redis.Get(ctx, "test").String()
	logx.Infof(ctx, "redis test: %s", t)

	logx.Err(ctx, errors.New("bbb"), "")

	web.JsonWithSuccess(ctx, "")
}
