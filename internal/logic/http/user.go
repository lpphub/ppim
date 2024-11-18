package http

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/lpphub/golib/render"
	"github.com/lpphub/golib/zlog"
	"go.mongodb.org/mongo-driver/mongo"
	"ppim/internal/logic/http/srv"
	"ppim/internal/logic/types"
	"ppim/pkg/errs"
)

type UserHandler struct {
	Srv *srv.UserSrv
}

func (h UserHandler) GetOne(ctx *gin.Context) {
	uid := ctx.Query("uid")
	if uid == "" {
		render.JsonWithError(ctx, errs.ErrInvalidParam)
		return
	}
	if u, err := h.Srv.GetOne(ctx, uid); err != nil {
		zlog.Error(ctx, err.Error())
		if errors.Is(err, mongo.ErrNoDocuments) {
			render.JsonWithError(ctx, errs.ErrRecordNotFound)
		} else {
			render.JsonWithError(ctx, errs.ErrServerInternal)
		}
	} else {
		render.JsonWithSuccess(ctx, u)
	}
}

func (h UserHandler) Register(ctx *gin.Context) {
	var req types.UserDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		zlog.Error(ctx, err.Error())
		render.JsonWithError(ctx, errs.ErrInvalidParam)
		return
	}

	if err := h.Srv.Register(ctx, req); err != nil {
		zlog.Error(ctx, err.Error())
		render.JsonWithError(ctx, errs.ErrServerInternal)
	} else {
		render.JsonWithSuccess(ctx, nil)
	}
}
