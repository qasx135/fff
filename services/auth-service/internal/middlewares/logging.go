package middlewares

import (
	"auth-service/internal/logger"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"gopkg.in/Graylog2/go-gelf.v2/gelf"
)

func GelfLoggerMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		m := &gelf.Message{
			Version: "1.0",
			Level:   gelf.LOG_INFO,
			Short:   fmt.Sprintf("%s %s", ctx.Request.Method, ctx.Request.URL.Path),
			Host:    "auth-service",
			Extra: map[string]any{
				"host":        "auth-service",
				"application": "auth-service",
				"method":      ctx.Request.Method,
				"path":        ctx.Request.URL.Path,
				"ip":          ctx.ClientIP(),
				"ua":          ctx.Request.UserAgent(),
			},
		}
		if err := logger.WriteMessage(m); err != nil {
			log.Printf("Ошибка отправки в грэйлог: %v", err)
		}

		ctx.Next()
	}
}
