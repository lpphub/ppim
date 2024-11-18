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
	authed := user.DID == req.Did && user.Token == req.Token
	resp := &logic.AuthResp{
		Ok: authed,
	}
	return resp, nil
}
