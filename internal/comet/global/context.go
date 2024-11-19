package global

import (
	"github.com/lpphub/golib/env"
	"path/filepath"
	"ppim/internal/comet/conf"
)

var (
	Conf conf.Config
)

func Init() {
	// 加载配置文件
	configFile := filepath.Join("config", "comet.yml")
	env.LoadConf(configFile, &Conf)
}
