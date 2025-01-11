package svc

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	"github.com/lpphub/golib/env"
	"ppim/internal/logic/global"
	"testing"
)

func TestConversationSrv_IndexConversation(t *testing.T) {
	env.SetRootPath("../../..")
	global.InitCtx()

	srv := ConversationSrv{}

	d, err := srv.ListByUID(context.TODO(), "456", 1736524291451, 1)
	if err != nil {
		t.Error(err)
		return
	}
	ds, _ := jsoniter.MarshalToString(d)
	t.Logf("%s", ds)
}
