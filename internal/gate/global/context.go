package global

import (
	"github.com/lpphub/golib/env"
	"github.com/lpphub/golib/logger"
	"path/filepath"
	"ppim/internal/gate/conf"
)

var (
	Conf conf.Config
)

func Init() {
	// 加载配置文件
	configFile := filepath.Join("config", "gate.yml")
	env.LoadConf(configFile, &Conf)

	// 配置日志
	logger.Setup()
}
