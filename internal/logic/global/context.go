package global

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/lpphub/golib/env"
	"github.com/lpphub/golib/zlog"
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
	Conf      conf.Config
	Mongo     *mongo.Database
	Redis     *redis.Client
	Snowflake *snowflake.Node
)

func PreInit() {
	log.Printf("Current Active RUN_ENV: %s \n", env.GetRunEnv())
	env.SetAppName("beacon")

	// 加载配置文件
	configFile := filepath.Join("config", fmt.Sprintf("logic_%s.yml", env.GetRunEnv()))
	env.LoadConf(configFile, &Conf)

	// 日志配置
	zlog.InitLog(zlog.WithLogPath(Conf.Log.Path))
}

func InitGlobalCtx() {
	PreInit()

	initDb()
	initCache()
	Snowflake, _ = snowflake.NewNode(1)
}

func Clear() {
	zlog.Close()
}
