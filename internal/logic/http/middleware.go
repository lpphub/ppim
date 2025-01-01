package http

import (
	"github.com/gin-gonic/gin"
)

func authorize() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
	}
}
