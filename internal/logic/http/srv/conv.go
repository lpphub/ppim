package srv

import (
	"github.com/gin-gonic/gin"
	"ppim/internal/logic/svc"
	"ppim/internal/logic/types"
)

type ConvSrv struct{}

func NewConvSrv() *ConvSrv {
	return &ConvSrv{}
}

func (srv *ConvSrv) RecentList(ctx *gin.Context, uid string) ([]*types.ConvRecentDTO, error) {
	return svc.Hints().Conv.GetRecentByUID(ctx, uid)
}

func (srv *ConvSrv) SetPin(ctx *gin.Context, req types.ConvOpVO) error {
	return svc.Hints().Conv.SetPin(ctx, req.UID, req.ConversationID, req.Pin)
}

func (srv *ConvSrv) SetMute(ctx *gin.Context, req types.ConvOpVO) error {
	return svc.Hints().Conv.SetMute(ctx, req.UID, req.ConversationID, req.Mute)
}
