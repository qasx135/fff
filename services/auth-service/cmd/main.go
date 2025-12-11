package main

import (
	"auth-service/internal/databases"
	"auth-service/internal/handlers"
	"auth-service/internal/logger"
	"auth-service/internal/middlewares"
	"auth-service/internal/telemetry"
	"auth-service/internal/utils"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func main() {
	config := utils.Load()
	logger.InitGelfLogger(config.GelfEndpoint)
	defer logger.Close()
	databases.InitRedis(config.RedisAddr)
	secretManager := &utils.JwtManager{}
	secretManager.InitSecret(config.JWTSecret, config.AccessTokenExp, config.RefreshTokenExp)
	shutdown := telemetry.InitTracer("auth-service", config.JaegerEndpoint)
	defer shutdown()

	authHandler := &handlers.AuthHandler{SignManager: secretManager}

	r := gin.Default()
	r.Use(middlewares.GelfLoggerMiddleware())

	r.POST("/login", authHandler.Login)
	r.POST("/refresh", authHandler.Refresh)

	r.GET("/healthz", handlers.Healthz)
	r.GET("/ready", handlers.Ready)
	r.GET("/startup", handlers.Startup)
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	r.Use(otelgin.Middleware("auth-service"))

	log.Println("Auth service запущен на :8080")
	if err := r.Run(config.ServerPort); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
