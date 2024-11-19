package global

import (
	"github.com/bwmarrin/snowflake"
	"github.com/lpphub/golib/env"
	"path/filepath"
	"ppim/internal/comet/conf"
)

var (
	Conf      conf.Config
	Snowflake *snowflake.Node
)

func Init() {
	// 加载配置文件
	configFile := filepath.Join("config", "comet.yml")
	env.LoadConf(configFile, &Conf)

	Snowflake, _ = snowflake.NewNode(1)
}
