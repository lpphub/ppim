package logic

import (
	"github.com/lpphub/golib/ware"
	"github.com/lpphub/golib/zlog"
	"net/http"
	"ppim/internal/logic/global"
	"ppim/internal/logic/router"
)

func HttpServe() {
	global.InitGlobalCtx()
	defer global.Clear()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router.SetupRouters(),
	}
	zlog.Infof(nil, "Listening and serving HTTP on %s", srv.Addr)
	ware.ListenAndServe(srv)
}
