package ext

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lpphub/golib/zlog"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"time"
)

type RedisLogHook struct{}

func (RedisLogHook) DialHook(next redis.DialHook) redis.DialHook {
	return next
}

func (RedisLogHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		start := time.Now()
		err := next(ctx, cmd)
		end := time.Now()

		msg := "redis do success"
		if err != nil {
			msg = "redis do error: " + err.Error()
		}
		if c, ok := ctx.(*gin.Context); ok && c != nil {
			fields := []zap.Field{
				zap.String("logId", zlog.GetLogId(c)),
				zap.String("url", c.Request.URL.Path),
				zap.String("command", fmt.Sprintf("%s", joinArgs(1024, cmd.Args()))),
				zap.Float64("cost", float64(end.Sub(start).Nanoseconds()/1e4)/100.0),
			}
			zlog.ZapLogger.Info(msg, fields...)
		}
		return err
	}
}

func (RedisLogHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return next
}

func joinArgs(showByte int, args ...interface{}) string {
	if showByte < 0 {
		return ""
	}
	var sumLen int
	var argStr string
	for _, v := range args {
		tmp := fmt.Sprintf("%v", v) + " "
		argStr += tmp
		sumLen += len(tmp)

		if showByte > 0 && sumLen >= showByte {
			break
		}
	}
	if showByte > 0 && sumLen > showByte {
		argStr = argStr[:showByte] + "..."
	}
	return argStr
}
