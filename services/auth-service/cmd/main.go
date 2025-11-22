package main

import (
	"auth-service/internal/databases"
	"auth-service/internal/handlers"
	"auth-service/internal/logger"
	"auth-service/internal/middlewares"
	"auth-service/internal/utils"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	logger.InitGelfLogger()
	databases.InitRedis()
	utils.InitSecret()

	r := gin.Default()
	r.Use(middlewares.GelfLoggerMiddleware())

	r.POST("/login", handlers.Login)
	r.POST("/refresh", handlers.Refresh)

	r.GET("/healthz", handlers.Healthz)
	r.GET("/ready", handlers.Ready)
	r.GET("/startup", handlers.Startup)

	log.Println("Auth service запущен на :8080")
	r.Run(":8080")
}
