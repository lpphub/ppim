package logic

import (
	//"net/http"
	"ppim/internal/logic/global"
	"ppim/internal/logic/http"
)

func Serve() {
	global.InitGlobalCtx()
	defer global.Clear()

	api := http.NewApiServer()
	api.Start()
}
