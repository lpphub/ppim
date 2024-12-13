package http

import (
	"github.com/gin-gonic/gin"
)

func TestA() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
	}
}
