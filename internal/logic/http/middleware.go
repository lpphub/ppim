package http

import (
	"github.com/gin-gonic/gin"
)

func Test() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
	}
}
