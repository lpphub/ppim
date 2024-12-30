package rpc

import (
	"context"
	"github.com/valyala/fastrand"
	"ppim/internal/chatlib"
	"ppim/pkg/util"
)

// 自定义选择器：同一会话路由到同一logic处理器
type customSelector struct {
	servers []string
}

func (s *customSelector) Select(_ context.Context, _, serviceMethod string, args interface{}) string {
	ss := s.servers
	sl := len(ss)
	if sl == 0 {
		return ""
	}

	// 根据会话ID路由
	if serviceMethod == methodSendMsg {
		cid := args.(*chatlib.MessageReq).ConversationID
		return ss[(int(util.CRC32(cid))%sl+sl)%sl]
	}

	// 随机选择一个
	i := fastrand.Uint32n(uint32(len(ss)))
	return ss[i]
}

func (s *customSelector) UpdateServer(servers map[string]string) {
	ss := make([]string, 0, len(servers))
	for k := range servers {
		ss = append(ss, k)
	}
	s.servers = ss
}
