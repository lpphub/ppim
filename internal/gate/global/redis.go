package global

import (
	"github.com/redis/go-redis/v9"
	"ppim/pkg/ext"
	"time"
)

func initRedis() {
	Redis = redis.NewClient(&redis.Options{
		Addr:            Conf.Redis.Addr,
		Password:        Conf.Redis.Password,
		DB:              Conf.Redis.DB,
		MinIdleConns:    2,
		MaxActiveConns:  10,
		ConnMaxLifetime: 5 * time.Minute,
	})
	Redis.AddHook(ext.RedisLogHook{})
}
