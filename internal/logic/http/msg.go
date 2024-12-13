package http

import (
	"github.com/gin-gonic/gin"
	"github.com/lpphub/golib/logger/logx"
	"github.com/lpphub/golib/web"
	"ppim/internal/logic/http/srv"
	"ppim/pkg/errs"
)

type MsgHandler struct {
	conv *srv.ConvSrv
}

func (h *MsgHandler) RecentConvList(ctx *gin.Context) {
	uid := ctx.Query("uid")
	if uid == "" {
		web.JsonWithError(ctx, errs.ErrInvalidParam)
		return
	}

	if list, err := h.conv.RecentList(ctx, uid); err != nil {
		logx.Err(ctx, err, "")
		web.JsonWithError(ctx, errs.ErrRecordNotFound)
	} else {
		web.JsonWithSuccess(ctx, list)
	}
}
