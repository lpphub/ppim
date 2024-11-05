package rpc

import (
	"context"
	"ppim/api/logic"
)

type pingServer struct {
	logic.UnimplementedPingServer
}

func (s *pingServer) Ping(ctx context.Context, req *logic.PingReq) (*logic.PongResp, error) {
	return &logic.PongResp{
		Msg: "pong haha",
	}, nil

}
