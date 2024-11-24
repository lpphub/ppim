package http

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/lpphub/golib/logger/glog"
	"github.com/lpphub/golib/web"
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
		web.JsonWithError(ctx, errs.ErrInvalidParam)
		return
	}
	if u, err := h.Srv.GetOne(ctx, uid); err != nil {
		glog.Error(ctx, err.Error())
		if errors.Is(err, mongo.ErrNoDocuments) {
			web.JsonWithError(ctx, errs.ErrRecordNotFound)
		} else {
			web.JsonWithError(ctx, errs.ErrServerInternal)
		}
	} else {
		web.JsonWithSuccess(ctx, u)
	}
}

func (h UserHandler) Register(ctx *gin.Context) {
	var req types.UserDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		glog.Error(ctx, err.Error())
		web.JsonWithError(ctx, errs.ErrInvalidParam)
		return
	}

	if err := h.Srv.Register(ctx, req); err != nil {
		glog.Error(ctx, err.Error())
		web.JsonWithError(ctx, errs.ErrServerInternal)
	} else {
		web.JsonWithSuccess(ctx, nil)
	}
}
