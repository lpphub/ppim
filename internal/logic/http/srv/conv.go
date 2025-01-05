package srv

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"ppim/internal/logic/service"
	"ppim/internal/logic/store"
	"ppim/internal/logic/types"
)

type ConvSrv struct{}

func NewConvSrv() *ConvSrv {
	return &ConvSrv{}
}

func (srv *ConvSrv) RecentList(ctx *gin.Context, uid string) (list []*types.RecentConvVO, err error) {
	list, err = service.Hints().Conv.CacheQueryRecent(ctx, uid)
	if err == nil {
		return
	}

	data, err := new(store.Conversation).ListRecent(ctx, uid)
	if err != nil {
		return nil, err
	}

	msgIds := make([]string, 0, len(data))
	for _, d := range data {
		list = append(list, &types.RecentConvVO{
			ConversationID:   d.ConversationID,
			ConversationType: d.ConversationType,
			Mute:             d.Mute,
			Pin:              d.Pin,
			FromUID:          d.FromUID,
			LastMsgID:        d.LastMsgId,
			LastMsgSeq:       d.LastMsgSeq,
			Version:          d.CreatedAt.UnixMilli(),
		})

		msgIds = append(msgIds, d.LastMsgId)
	}

	msgList, err := new(store.Message).ListByMsgIds(ctx, msgIds)
	if err != nil {
		return nil, err
	}

	msgMap := make(map[string]*types.MessageDTO, len(msgList))
	for _, m := range msgList {
		var md types.MessageDTO
		_ = copier.Copy(&m, md)
		msgMap[m.MsgID] = &md
	}

	for _, vo := range list {
		if md, ok := msgMap[vo.LastMsgID]; ok {
			vo.LastMsg = md
		}
	}
	return
}

func (srv *ConvSrv) SetPin(ctx *gin.Context, req types.ConvOpDTO) error {
	if err := service.Hints().Conv.CacheSetPin(ctx, req.UID, req.ConversationID, req.Pin); err != nil {
		return err
	}
	return new(store.Conversation).UpdatePin(ctx, req.UID, req.ConversationID, req.Pin)
}

func (srv *ConvSrv) SetMute(ctx *gin.Context, req types.ConvOpDTO) error {
	if err := service.Hints().Conv.CacheSetMute(ctx, req.UID, req.ConversationID, req.Mute); err != nil {
		return err
	}
	return new(store.Conversation).UpdateMute(ctx, req.UID, req.ConversationID, req.Mute)
}
