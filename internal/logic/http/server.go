package http

import (
	"github.com/gin-gonic/gin"
	"github.com/lpphub/golib/render"
	"github.com/lpphub/golib/ware"
	"net/http"
	"ppim/pkg/errs"
	"ppim/pkg/ext"
)

type ApiServer struct {
	engine *gin.Engine
}

func NewApiServer() *ApiServer {
	r := gin.New()
	ware.Bootstraps(r, ware.BootstrapConf{
		Cors: true,
		AccessLog: ware.AccessLogConfig{
			Enable:    true,
			SkipPaths: []string{"/metrics"},
		},
		CustomRecovery: func(ctx *gin.Context, err any) {
			render.JsonWithError(ctx, errs.ErrServerInternal)
		},
	})

	ext.RegisterMetrics(r)
	registerRoute(r)

	return &ApiServer{
		engine: r,
	}
}

func (s *ApiServer) Start() {
	srv := &http.Server{
		Addr:    ":8080",
		Handler: s.engine.Handler(),
	}
	ware.ListenAndServe(srv)
}

func (s *ApiServer) Stop() {
}
