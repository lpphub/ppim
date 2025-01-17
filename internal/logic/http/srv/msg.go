package srv

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"ppim/internal/logic/store"
	"ppim/internal/logic/types"
)

type MsgSrv struct{}

func NewMsgSrv() *MsgSrv {
	return &MsgSrv{}
}

func (s *MsgSrv) ListConvMsg(ctx *gin.Context, req types.MessageQueryVO) ([]types.MessageDTO, error) {
	list, err := new(store.Message).ListByConvSeq(ctx, req.ConversationID, req.StartSeq, req.Limit)
	if err != nil {
		return nil, err
	}

	voList := make([]types.MessageDTO, 0, len(list))
	for _, v := range list {
		var vo types.MessageDTO
		_ = copier.Copy(&vo, v)
		vo.SendTime = v.SendTime.UnixMilli()
		vo.CreatedAt = v.CreatedAt.UnixMilli()
		vo.UpdatedAt = v.UpdatedAt.UnixMilli()
		voList = append(voList, vo)
	}
	return voList, nil
}
