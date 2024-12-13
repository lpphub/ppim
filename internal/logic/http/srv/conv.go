package srv

import (
	"github.com/gin-gonic/gin"
	"ppim/internal/logic/store"
	"ppim/internal/logic/types"
)

type ConvSrv struct{}

func NewConvSrv() *ConvSrv {
	return &ConvSrv{}
}

func (srv *ConvSrv) RecentList(ctx *gin.Context, uid string) (list []*types.RecentConvVO, err error) {
	//ids, _ := global.Redis.ZRevRange(ctx, fmt.Sprintf(service.CacheConvRecent, uid), 0, 200).Result()
	//logx.Infof(ctx, "recent conv ids=%v", ids)

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
			FromUid:          d.FromID,
			LastMsgID:        d.LastMsgId,
		})

		msgIds = append(msgIds, d.LastMsgId)
	}

	return
}
