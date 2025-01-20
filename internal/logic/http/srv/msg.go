package srv

import (
	"github.com/gin-gonic/gin"
	"ppim/internal/logic/svc"
	"ppim/internal/logic/types"
)

type MsgSrv struct{}

func NewMsgSrv() *MsgSrv {
	return &MsgSrv{}
}

func (s *MsgSrv) ListConvMsg(ctx *gin.Context, req types.MessageQueryVO) ([]types.MessageDTO, error) {
	return svc.Hints().Msg.PullUpOrDown(ctx, req.ConversationID, req.StartSeq, req.Limit)
}

func (s *MsgSrv) Revoke(ctx *gin.Context, req types.MsgOpVO) error {
	return svc.Hints().Msg.Revoke(ctx, req.MsgID, req.ConversationID)
}

func (s *MsgSrv) Delete(ctx *gin.Context, req types.MsgOpVO) error {
	return svc.Hints().Msg.Delete(ctx, req.MsgID, req.ConversationID)
}
