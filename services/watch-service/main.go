package main

import (
	"fmt"
	"github.com/Graylog2/go-gelf/gelf"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

type WatchEvent struct {
	gorm.Model
	UserID  string `json:"user_id"`
	AnimeID uint   `json:"anime_id"`
	Episode int    `json:"episode"`
}

var (
	db         *gorm.DB
	gelfWriter *gelf.Writer
)

func initLogger() {
	var err error
	gelfWriter, err = gelf.NewWriter("graylog:12201")
	if err != nil {
		log.Fatal("Не удалось подключиться к Graylog:", err)
	}
	gelfWriter.Facility = "watch-service"
}

func logRequest(c *gin.Context) {
	msg := gelf.Message{
		Version: "1.1",
		Level:   gelf.LOG_INFO,
		Short:   fmt.Sprintf("%s %s", c.Request.Method, c.Request.URL.Path),
		Extra: map[string]interface{}{
			"method":  c.Request.Method,
			"url":     c.Request.URL.Path,
			"ip":      c.ClientIP(),
			"service": "watch-service",
		},
	}
	if err := gelfWriter.WriteMessage(&msg); err != nil {
		log.Printf("Ошибка отправки в Graylog: %v", err)
	}
	c.Next()
}

func initDB() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "host=watch-db-postgresql user=watch_user password=watch_pass dbname=watch_db port=5432 sslmode=disable"
	}
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Не удалось подключиться к БД:", err)
	}
	db.AutoMigrate(&WatchEvent{})
}

func healthz(c *gin.Context) { c.Status(200) }
func ready(c *gin.Context)   { c.Status(200) }
func startup(c *gin.Context) { c.Status(200) }

func getHistory(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(400, gin.H{"error": "user_id required"})
		return
	}

	var events []WatchEvent
	db.Where("user_id = ?", userID).Find(&events)
	c.JSON(200, events)
}

func watchHandler(c *gin.Context) {
	var event WatchEvent
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if event.UserID == "" || event.AnimeID == 0 {
		c.JSON(400, gin.H{"error": "user_id and anime_id required"})
		return
	}

	db.Create(&event)
	c.JSON(201, event)
}

func main() {
	initLogger()
	defer gelfWriter.Close()

	initDB()

	r := gin.Default()
	r.Use(logRequest)

	r.POST("/watch", watchHandler)
	r.GET("/history/:user_id", getHistory)

	r.GET("/healthz", healthz)
	r.GET("/ready", ready)
	r.GET("/startup", startup)

	log.Println("Watch service запущен на :8080")
	r.Run(":8080")
}
