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

func (h ChatHandler) ConvRecentList(ctx *gin.Context) {
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

func (h ChatHandler) ConvListMsg(ctx *gin.Context) {
	var req types.ConvMessageVO
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

func (h ChatHandler) ConvPin(ctx *gin.Context) {
	var req types.ConvOpVO
	if err := ctx.ShouldBind(&req); err != nil {
		web.JsonWithError(ctx, errs.ErrInvalidParam)
		return
	}
	if err := h.conv.SetPin(ctx, req); err != nil {
		logx.Err(ctx, err, "")
		web.JsonWithError(ctx, errs.ErrServerInternal)
		return
	}
	web.JsonWithSuccess(ctx, "ok")
}

func (h ChatHandler) ConvMute(ctx *gin.Context) {
	var req types.ConvOpVO
	if err := ctx.ShouldBind(&req); err != nil {
		web.JsonWithError(ctx, errs.ErrInvalidParam)
		return
	}
	if err := h.conv.SetMute(ctx, req); err != nil {
		logx.Err(ctx, err, "")
		web.JsonWithError(ctx, errs.ErrServerInternal)
		return
	}
	web.JsonWithSuccess(ctx, "ok")
}

func (h ChatHandler) ConvSetUnread(ctx *gin.Context) {
	// todo
	web.JsonWithSuccess(ctx, "ok")
}

func (h ChatHandler) ConvDel(ctx *gin.Context) {
	// todo
	web.JsonWithSuccess(ctx, "ok")
}

func (h ChatHandler) MsgDel(ctx *gin.Context) {
	// todo
	web.JsonWithSuccess(ctx, "ok")
}

func (h ChatHandler) MsgRevoke(ctx *gin.Context) {
	// todo
	web.JsonWithSuccess(ctx, "ok")
}

func (h ChatHandler) MsgRead(ctx *gin.Context) {
	// todo
	web.JsonWithSuccess(ctx, "ok")
}
