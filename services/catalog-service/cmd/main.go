package main

import (
	"catalog-service/internal/databases"
	"catalog-service/internal/handlers"
	"catalog-service/internal/kafka"
	"catalog-service/internal/logger"
	"catalog-service/internal/middlewares"
	"catalog-service/internal/telemetry"
	"catalog-service/internal/utils"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func main() {
	config := utils.Load()
	logger.InitGelfLogger(config.GelfEndpoint)
	defer logger.Close() // Важно: закрыть соединение при завершении
	shutdown := telemetry.InitTracer("catalog-service", config.JaegerEndpoint)
	defer shutdown()

	reader := kafka.InitReader(config.BrokerHost, "anime-watch", "anime-watch")

	defer reader.Close()

	jwtManager := &utils.JwtManager{}
	jwtManager.InitSecret(config.JWTSecret)

	databases.InitDB(config.PostgresqlDSN)

	r := gin.Default()
	r.Use(otelgin.Middleware("catalog-service"))
	r.Use(middlewares.GelfLoggerMiddleware())
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	animeHandler := &handlers.AnimeHandler{}
	animeHandler.SetSources(databases.DB, reader)

	go animeHandler.HandleKafkaEvents()

	api := r.Group("api")
	api.Use(middlewares.AuthMiddleware(jwtManager))

	{
		api.GET("/anime/:id", animeHandler.GetAnime)
		api.GET("/anime", animeHandler.GetAnimes)
	}

	r.GET("/healthz", handlers.Healthz)
	r.GET("/ready", handlers.Ready)
	r.GET("/startup", handlers.Startup)

	log.Println("Catalog service запущен на :8080")
	r.Run(config.ServerPort)
}
