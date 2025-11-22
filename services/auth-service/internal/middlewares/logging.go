package middlewares

import (
	"auth-service/internal/logger"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/Graylog2/go-gelf.v2/gelf"
)

func GelfLoggerMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		startTime := time.Now()
		ctx.Next()
		endTime := time.Now()
		latencyTime := endTime.Sub(startTime).Milliseconds()
		reqMethod := ctx.Request.Method
		reqUri := ctx.Request.RequestURI
		statusCode := ctx.Writer.Status()
		clientIP := ctx.ClientIP()
		logger.WriteMessage(
			&gelf.Message{
				Host:     "API",
				Short:    reqMethod,
				TimeUnix: float64(endTime.Unix()),
				Extra: map[string]any{
					"_reqUri":   reqUri,
					"_clientIP": clientIP,
					"_status":   statusCode,
					"_latency":  latencyTime,
				},
			},
		)

		ctx.Next()
	}
}