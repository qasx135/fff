package handlers

import (
	"log"
	"time"
	internal_kafka "watch-service/internal/kafka"
	"watch-service/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"gorm.io/gorm"
)

type WatchHandler struct {
	db *gorm.DB
}

func (h *WatchHandler) SetDB(DB *gorm.DB) {
	h.db = DB
}

type WatchEventDTO struct {
	UserID  string    `json:"user_id"`
	AnimeID uuid.UUID `json:"anime_id"`
	Episode int       `json:"episode"`
}

func (h *WatchHandler) GetHistory(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(400, gin.H{"error": "user_id required"})
		return
	}

	var events []models.WatchEvent
	h.db.Where("user_id = ?", userID).Find(&events)
	c.JSON(200, events)
}

func (h *WatchHandler) WatchHandler(c *gin.Context) {
	var dto WatchEventDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if dto.UserID == "" || dto.AnimeID == uuid.Nil {
		c.JSON(400, gin.H{"error": "user_id and anime_id required"})
		return
	}

	event := models.WatchEvent{
		UserID:  dto.UserID,
		AnimeID: dto.AnimeID,
		Episode: dto.Episode,
	}

	h.db.Create(&event)

	if err := internal_kafka.WriteMessage(kafka.Message{
		Key:   []byte(event.AnimeID.String()),
		Value: []byte(event.UpdatedAt.Format(time.RFC3339)),
	}); err != nil {
		log.Printf("Error while sending message: %s", err)
	}

	c.JSON(201, event)
}
