package main

import (
	"log"
	"watch-service/internal/databases"
	"watch-service/internal/handlers"
	"watch-service/internal/logger"
	"watch-service/internal/middlewares"
	"watch-service/internal/telemetry"
	"watch-service/internal/utils"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	
)

func main() {
	config := utils.Load()
	jwtManager := &utils.JwtManager{}
	jwtManager.InitSecret(config.JWTSecret)
	logger.InitGelfLogger(config.GelfEndpoint)
	defer logger.Close()
	shutdown := telemetry.InitTracer("watch-service", config.JaegerEndpoint)
	defer shutdown()

	databases.InitDB(config.PostgresqlDSN)
	watchHandler := &handlers.WatchHandler{}
	watchHandler.SetDB(databases.DB)

	r := gin.Default()
	r.Use(otelgin.Middleware("watch-service"))
	r.Use(middlewares.GelfLoggerMiddleware())
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	api := r.Group("/api")
	api.Use(middlewares.AuthMiddleware(jwtManager))
	{
		api.POST("/watch", watchHandler.WatchHandler)
		api.GET("/history/:user_id", watchHandler.GetHistory)
	}

	r.GET("/healthz", handlers.Healthz)
	r.GET("/ready", handlers.Ready)
	r.GET("/startup", handlers.Startup)

	log.Println("Watch service запущен на :8080")
	r.Run(config.ServerPort)
}
