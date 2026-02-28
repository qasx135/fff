package middlewares

import (
	"fmt"
	"log"
	"watch-service/internal/logger"

	"github.com/gin-gonic/gin"
	"gopkg.in/Graylog2/go-gelf.v2/gelf"
)

func GelfLoggerMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if err := logger.WriteMessage(
			&gelf.Message{
				Version: "1.0",
				Level:   gelf.LOG_INFO,
				Short:   fmt.Sprintf("%s %s", ctx.Request.Method, ctx.Request.URL.Path),
				Host:    "watch-service",
				Extra: map[string]any{
					"host":        "watch-service",
					"application": "watch-service",
					"method":      ctx.Request.Method,
					"path":        ctx.Request.URL.Path,
					"ip":          ctx.ClientIP(),
					"ua":          ctx.Request.UserAgent(),
				},
			},
		); err != nil {
			log.Printf("Ошибка отправки в грэйлог: %v", err)
		}

		ctx.Next()
	}
}
