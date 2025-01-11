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
	attr := types.ConvAttributeDTO{
		UID:            req.UID,
		ConversationID: req.ConversationID,
		Attribute:      svc.ConvFieldPin,
		Pin:            req.Pin,
	}
	err := svc.Hints().Conv.SetAttribute(ctx, attr)
	if err != nil {
		return err
	}
	// todo 多端同步
	return nil
}

func (srv *ConvSrv) SetMute(ctx *gin.Context, req types.ConvOpVO) error {
	attr := types.ConvAttributeDTO{
		UID:            req.UID,
		ConversationID: req.ConversationID,
		Attribute:      svc.ConvFieldMute,
		Mute:           req.Mute,
	}
	err := svc.Hints().Conv.SetAttribute(ctx, attr)
	if err != nil {
		return err
	}
	// todo 多端同步
	return nil
}

func (srv *ConvSrv) SetUnreadCount(ctx *gin.Context, req types.ConvOpVO) error {
	attr := types.ConvAttributeDTO{
		UID:            req.UID,
		ConversationID: req.ConversationID,
		Attribute:      svc.ConvFieldUnreadCount,
		UnreadCount:    req.UnreadCount,
	}
	err := svc.Hints().Conv.SetAttribute(ctx, attr)
	if err != nil {
		return err
	}
	// todo 多端同步
	return nil
}

func (srv *ConvSrv) SetDel(ctx *gin.Context, req types.ConvOpVO) error {
	attr := types.ConvAttributeDTO{
		UID:            req.UID,
		ConversationID: req.ConversationID,
		Attribute:      svc.ConvFieldDeleted,
		Deleted:        req.Deleted,
	}
	err := svc.Hints().Conv.SetAttribute(ctx, attr)
	if err != nil {
		return err
	}
	// todo 多端同步
	return nil
}
