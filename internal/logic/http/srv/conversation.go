package srv

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lpphub/golib/logger/logx"
	"ppim/internal/logic/global"
	"ppim/internal/logic/service"
	"ppim/internal/logic/store"
	"ppim/internal/logic/types"
)

type ConversationSrv struct{}

func NewConversationSrv() *ConversationSrv {
	return &ConversationSrv{}
}

func (srv *ConversationSrv) RecentList(ctx *gin.Context, uid string) (list []*types.RecentConvVO, err error) {
	ids, _ := global.Redis.ZRevRange(ctx, fmt.Sprintf(service.CacheConvRecent, uid), 0, 200).Result()
	logx.Infof(ctx, "recent conv ids=%v", ids)

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
