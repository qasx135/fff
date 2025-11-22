package main

import (
	"fmt"
	"log"
	_ "net"
	"os"
	_ "time"

	"github.com/Graylog2/go-gelf/gelf"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db         *gorm.DB
	gelfWriter *gelf.Writer
)

type Anime struct {
	gorm.Model
	Title       string
	Description string
	Genre       string
	Episodes    int
}

func initLogger() {
	var err error
	// Подключаемся к Graylog по UDP на порту 12201
	gelfWriter, err = gelf.NewWriter("graylog:12201")
	if err != nil {
		log.Fatal("Не удалось подключиться к Graylog:", err)
	}
	// Опционально: задаём метаданные
	gelfWriter.Facility = "catalog-service"
}

func logRequest(c *gin.Context) {
	// Собираем структурированное GELF-сообщение
	msg := gelf.Message{
		Version: "1.1",
		Level:   gelf.LOG_INFO,
		Short:   fmt.Sprintf("%s %s", c.Request.Method, c.Request.URL.Path),
		Extra: map[string]interface{}{
			"method":  c.Request.Method,
			"url":     c.Request.URL.Path,
			"ip":      c.ClientIP(),
			"service": "catalog-service",
		},
	}

	// Отправляем в Graylog
	if err := gelfWriter.WriteMessage(&msg); err != nil {
		// Логируем ошибку локально, чтобы не падать
		log.Printf("Ошибка отправки в Graylog: %v", err)
	}

	c.Next()
}

func initDB() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "host=catalog-db-postgresql user=catalog_user password=catalog_pass dbname=catalog_db port=5432 sslmode=disable"
	}
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Не удалось подключиться к БД:", err)
	}
	db.AutoMigrate(&Anime{})

	// Демо-данные
	db.FirstOrCreate(&Anime{}, Anime{
		Title:    "Cowboy Bebop",
		Genre:    "Sci-Fi",
		Episodes: 26,
	})
}

func healthz(c *gin.Context) { c.Status(200) }
func ready(c *gin.Context)   { c.Status(200) }
func startup(c *gin.Context) { c.Status(200) }

func getAnime(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		var list []Anime
		db.Find(&list)
		c.JSON(200, list)
		return
	}

	var anime Anime
	if db.First(&anime, id).Error != nil {
		c.JSON(404, gin.H{"error": "not found"})
		return
	}
	c.JSON(200, anime)
}

func main() {
	initLogger()
	defer gelfWriter.Close() // Важно: закрыть соединение при завершении

	initDB()

	r := gin.Default()
	r.Use(logRequest)

	r.GET("/anime/:id", getAnime)
	r.GET("/anime", getAnime)

	r.GET("/healthz", healthz)
	r.GET("/ready", ready)
	r.GET("/startup", startup)

	log.Println("Catalog service запущен на :8080")
	r.Run(":8080")
}
