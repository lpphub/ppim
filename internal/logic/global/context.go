package global

import (
	"github.com/lpphub/golib/env"
	"github.com/lpphub/golib/logger"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"path/filepath"
	"ppim/internal/logic/conf"
)

//type Context struct {
//	Conf      conf.Config
//	DB        *gorm.DB
//	Redis     *redis.Client
//	Snowflake *snowflake.Node
//}
//
//func NewGlobalContext(c conf.Config) *Context {
//	zlog.InitLog(zlog.WithLogPath(c.Log.Path))
//	ctx := &Context{
//		Conf: c,
//	}
//	ctx.Snowflake, _ = snowflake.NewNode(1)
//	return ctx
//}

var (
	Conf  conf.Config
	Mongo *mongo.Database
	Redis *redis.Client
)

func PreInit() {
	log.Printf("Current Active RUN_ENV: %s \n", env.GetRunEnv())
	env.SetAppName("beacon")

	// 加载配置文件
	configFile := filepath.Join("config", "logic.yml")
	env.LoadConf(configFile, &Conf)

	// 日志配置
	logger.Setup()
}

func InitGlobalCtx() {
	PreInit()

	initDb()
	initRedis()
}
