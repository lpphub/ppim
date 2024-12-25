package util

import "github.com/lpphub/golib/logger"

func WithRecover(f func()) {
	defer func() {
		if err := recover(); err != nil {
			logger.Log().Error().Msgf("goroutine error: %v", err)
		}
	}()

	f()
}
