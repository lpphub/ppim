package http

import (
	"github.com/gin-gonic/gin"
	"github.com/lpphub/golib/web"
	"net/http"
	"ppim/internal/logic/http/errs"
	"ppim/pkg/ext"
)

type ApiServer struct {
	addr   string
	engine *gin.Engine
}

func NewApiServer(addr string) *ApiServer {
	r := gin.New()
	bootstraps(r)
	registerRoutes(r)
	return &ApiServer{
		addr:   addr,
		engine: r,
	}
}

func bootstraps(r *gin.Engine) {
	web.Bootstraps(r, web.BootstrapConf{
		Cors: true,
		AccessLog: web.AccessLogConfig{
			Enable:    true,
			SkipPaths: []string{"/metrics"},
		},
		CustomRecovery: func(ctx *gin.Context, err any) {
			web.JsonWithError(ctx, errs.ErrServerInternal)
		},
	})

	ext.RegisterMetrics(r)
}

func (s *ApiServer) Start() {
	srv := &http.Server{
		Addr:    s.addr,
		Handler: s.engine.Handler(),
	}
	web.ListenAndServe(srv)
}

func (s *ApiServer) Stop() {
}
