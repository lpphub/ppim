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
		rpcsrv := rpc.NewRpcServer(global.Conf.Server.Rpc, global.Conf.Server.Etcd)
		rpcsrv.Start()
	}()

	api := http.NewApiServer(global.Conf.Server.Api)
	api.Start()
}
