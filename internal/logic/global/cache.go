package global

import (
	"github.com/redis/go-redis/v9"
	"ppim/pkg/ext"
)

func initCache() {
	Redis = redis.NewClient(&redis.Options{
		Addr:     Conf.Redis.Addr,
		Password: Conf.Redis.Password,
		DB:       Conf.Redis.DB,
	})
	Redis.AddHook(ext.RedisLogHook{})
}
