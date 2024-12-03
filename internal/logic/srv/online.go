package srv

import (
	"context"
	"fmt"
	"ppim/internal/logic/global"
	"ppim/internal/logic/types"
)

type OnlineSrv struct {
}

const (
	cacheRouteUid = "online:%s"
)

var OnSrv = &OnlineSrv{}

func (s *OnlineSrv) Register(ctx context.Context, ol *types.OnlineDTO) error {
	err := global.Redis.SAdd(ctx, s.getOnlineKey(ol.Uid), s.buildVal(ctx, ol)).Err()
	return err
}

func (s *OnlineSrv) UnRegister(ctx context.Context, ol *types.OnlineDTO) error {
	err := global.Redis.SAdd(ctx, s.getOnlineKey(ol.Uid), s.buildVal(ctx, ol)).Err()
	return err
}

func (s *OnlineSrv) getOnlineKey(uid string) string {
	return fmt.Sprintf(cacheRouteUid, uid)
}

func (s *OnlineSrv) buildVal(_ context.Context, ol *types.OnlineDTO) string {
	return fmt.Sprintf("%s_%s_%s_%s", ol.Uid, ol.Did, ol.Topic, ol.Ip)
}
