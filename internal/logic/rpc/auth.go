package rpc

import (
	"context"
	"ppim/api/logic"
	"ppim/internal/logic/model"
)

type authService struct {
	logic.UnimplementedAuthServer
}

func (s *authService) Auth(ctx context.Context, req *logic.AuthReq) (*logic.AuthResp, error) {
	user := new(model.User)
	if err := user.GetOne(ctx, req.Uid); err != nil {
		return nil, err
	}
	code := 0
	if user.DID != req.Did || user.Token != req.Token {
		code = 1 // 鉴权失败
	}
	resp := &logic.AuthResp{
		Code: uint32(code),
	}
	return resp, nil
}
