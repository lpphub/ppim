package router

import (
	"github.com/gin-gonic/gin"
	"github.com/lpphub/golib/render"
	"github.com/lpphub/golib/ware"
	"ppim/pkg/errs"
	"ppim/pkg/ext"
)

func SetupRouters() *gin.Engine {
	r := gin.New()
	ware.Bootstraps(r, ware.BootstrapConf{
		Cors: true,
		TraceLog: ware.TraceLogConfig{
			Enable:      true,
			IgnorePaths: []string{"/metrics"},
		},
		CustomRecovery: func(ctx *gin.Context, err any) {
			render.JsonWithError(ctx, errs.ErrServerInternal)
		},
	})

	ext.RegistryMetrics(r)

	registerRoute(r)
	return r
}
