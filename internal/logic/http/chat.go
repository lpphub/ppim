package http

import (
	"github.com/gin-gonic/gin"
	"github.com/lpphub/golib/logger/logx"
	"github.com/lpphub/golib/web"
	"ppim/internal/logic/http/errs"
	"ppim/internal/logic/http/srv"
	"ppim/internal/logic/types"
)

type ChatHandler struct {
	conv *srv.ConvSrv
	msg  *srv.MsgSrv
}

func (h *ChatHandler) RecentConvList(ctx *gin.Context) {
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

func (h *ChatHandler) ListConvMsg(ctx *gin.Context) {
	var req types.ConvMsgReq
	if err := ctx.ShouldBind(&req); err != nil {
		web.JsonWithError(ctx, errs.ErrInvalidParam)
		return
	}

	if list, err := h.msg.ListConvMsg(ctx, req); err != nil {
		logx.Err(ctx, err, "")
		web.JsonWithError(ctx, errs.ErrRecordNotFound)
	} else {
		web.JsonWithSuccess(ctx, list)
	}
}
