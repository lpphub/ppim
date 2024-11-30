package logic

import (
	//"net/http"
	"ppim/internal/logic/global"
	"ppim/internal/logic/http"
	"ppim/internal/logic/rpc"
)

func Serve() {
	global.InitGlobalCtx()

	go func() {
		rpcsrv := rpc.NewRpcServer()
		rpcsrv.Start()
	}()

	api := http.NewApiServer()
	api.Start()
}
