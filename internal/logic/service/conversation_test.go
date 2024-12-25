package service

import (
	"context"
	"github.com/lpphub/golib/env"
	"ppim/internal/logic/global"
	"ppim/internal/logic/types"
	"ppim/pkg/ext"
	"testing"
)

func TestConversationSrv_IndexConversation(t *testing.T) {
	env.SetRootPath("../../..")
	global.InitGlobalCtx()

	srv := ConversationSrv{
		segmentLock: ext.NewSegmentLock(3),
	}

	msg := &types.MessageDTO{
		MsgID:            "1235",
		MsgSeq:           2,
		MsgNo:            "cli02",
		MsgType:          1,
		Content:          "hello2",
		ConversationID:   "single|123@456",
		ConversationType: "single",
		FromUID:          "456",
		ToID:             "123",
	}
	srv.indexWithLock(context.TODO(), msg, "456")
}
