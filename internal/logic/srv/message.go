package srv

import (
	"context"
	"github.com/lpphub/golib/logger"
	"ppim/internal/chat"
	"ppim/internal/logic/global"
	"ppim/internal/logic/model"
	"ppim/internal/logic/types"
	"time"
)

type MessageSrv struct{}

var MsgSrv = &MessageSrv{}

func (s *MessageSrv) HandleMsg(ctx context.Context, msg *types.MessageDTO) error {
	// todo 1. 消息持久化
	mm := &model.Message{
		MsgId:            msg.MsgID,
		MsgType:          msg.MsgType,
		Content:          msg.Content,
		ConversationType: msg.ConversationType,
		FromID:           msg.FromID,
		ToID:             msg.ToID,
		CreatedAt:        time.Now(),
	}
	err := mm.Insert(ctx)
	if err != nil {
		return err
	}

	// todo 2. 消息推送（在线）
	if msg.ConversationType == chat.ConvSingle {
		online, _ := global.Redis.SMembers(ctx, OnSrv.getOnlineKey(msg.ToID)).Result()
		if len(online) > 0 {
			for _, ol := range online {
				logger.Infof(ctx, "push message to %s", ol)
			}
		}
	} else if msg.ConversationType == chat.ConvGroup {
		// 获取群组所有在线成员
		members, err := new(model.Group).ListMembers(ctx, msg.ToID)
		if err != nil {
			return err
		}
		logger.Infof(ctx, "push message to %v", members)
		// todo 选择在线群成员

	}

	// todo 3. 消息通知（离线）
	return nil
}
