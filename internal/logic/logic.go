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
		grpc := rpc.NewRpcServer()
		grpc.Start()
	}()

	api := http.NewApiServer()
	api.Start()
}
