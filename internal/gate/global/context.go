package global

import (
	"github.com/lpphub/golib/env"
	"github.com/lpphub/golib/logger"
	"github.com/redis/go-redis/v9"
	"path/filepath"
	"ppim/internal/gate/conf"
)

var (
	Conf  conf.Config
	Redis *redis.Client
)

func Init() {
	// 加载配置文件
	configFile := filepath.Join("config", "gate.yml")
	env.LoadConf(configFile, &Conf)

	// 配置日志
	logger.Setup()

	initRedis()
}
