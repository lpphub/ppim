package srv

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"ppim/internal/logic/svc"
	"ppim/internal/logic/types"
)

type ConvSrv struct{}

func NewConvSrv() *ConvSrv {
	return &ConvSrv{}
}

func (srv *ConvSrv) List(ctx *gin.Context, query types.ConvQueryVO) (*types.ConvListVO, error) {
	var (
		limit     = query.Limit + 1 // 多查一条用于判断是否有下一页
		startTime = cast.ToInt64(query.NextKey)
	)
	list, err := svc.Hints().Conv.IncrQuery(ctx, query.UID, startTime, int64(limit))
	if err != nil {
		return nil, err
	}
	listLen := len(list)
	if listLen == 0 {
		return nil, nil
	}

	nextKey := ""
	if listLen >= limit {
		nextKey = cast.ToString(list[listLen-1].Version) // version = last_time
		list = list[:listLen-1]
	}
	result := &types.ConvListVO{
		List:    list,
		NextKey: nextKey,
	}
	return result, nil
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
