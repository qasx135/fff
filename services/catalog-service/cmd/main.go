package main

import (
	"catalog-service/internal/databases"
	"catalog-service/internal/handlers"
	"catalog-service/internal/logger"
	"catalog-service/internal/middlewares"
	"catalog-service/internal/utils"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	config := utils.Load()
	logger.InitGelfLogger(config.GelfEndpoint)
	defer logger.Close() // Важно: закрыть соединение при завершении

	databases.InitDB(config.PostgresqlDSN)

	r := gin.Default()
	r.Use(middlewares.GelfLoggerMiddleware())

	animeHandler := &handlers.AnimeHandler{}

	r.GET("/anime/:id", animeHandler.GetAnime)
	r.GET("/anime", animeHandler.GetAnime)

	r.GET("/healthz", handlers.Healthz)
	r.GET("/ready", handlers.Ready)
	r.GET("/startup", handlers.Startup)

	log.Println("Catalog service запущен на :8080")
	r.Run(":8080")
}
