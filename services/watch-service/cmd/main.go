package main

import (
	"log"
	"watch-service/internal/databases"
	"watch-service/internal/handlers"
	"watch-service/internal/logger"
	"watch-service/internal/middlewares"
	"watch-service/internal/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	config := utils.Load()
	jwtManager := &utils.JwtManager{}
	jwtManager.InitSecret(config.JWTSecret)
	logger.InitGelfLogger(config.GelfEndpoint)
	defer logger.Close()

	databases.InitDB(config.PostgresqlDSN)
	watchHandler := &handlers.WatchHandler{}
	watchHandler.SetDB(databases.DB)

	r := gin.Default()
	r.Use(middlewares.GelfLoggerMiddleware())

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
	r.Run(":8080")
}
