package net

import (
	"github.com/lpphub/golib/logger"
	"testing"
)

func Test_newRetryManager(t *testing.T) {
	logger.Setup()

	svc := InitServerContext()
	svc.RetryManager.Start()

	select {}
}
