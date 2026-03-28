package logger

import (
	_ "embed"

	"github.com/gin-gonic/gin"
	"github.com/tuxounet/k2-sdk/types"
)

func Register(service types.IKernelService, server *gin.Engine) error {
	ingresLogger := service.GetLogger()
	server.Use(func(ctx *gin.Context) {
		ctx.Next()

		source := ctx.Request.Header.Get("X-Forwarded-For")
		if source == "" {
			source = ctx.ClientIP()
		}

		status := ctx.Writer.Status()
		switch {
		case status >= 500:
			ingresLogger.ErrorF("from %s - request %s %s %s %d", source, ctx.Request.Method, ctx.Request.URL.Path, ctx.Request.URL.Query(), ctx.Writer.Status())
		case status >= 400:
			ingresLogger.WarnF("from %s - request %s %s %s %d", source, ctx.Request.Method, ctx.Request.URL.Path, ctx.Request.URL.Query(), ctx.Writer.Status())
		default:
			ingresLogger.TraceF("from %s - request %s %s %s %d", source, ctx.Request.Method, ctx.Request.URL.Path, ctx.Request.URL.Query(), ctx.Writer.Status())
		}
	}, gin.Recovery())
	return nil
}
