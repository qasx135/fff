package handlers

import (
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
	var event models.WatchEvent
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if event.UserID == "" || event.AnimeID != uuid.Nil {
		c.JSON(400, gin.H{"error": "user_id and anime_id required"})
		return
	}

	h.db.Create(&event)

	internal_kafka.WriteMessage(&kafka.Message{
		Key:   []byte(event.AnimeID.String()),
		Value: []byte(event.UpdatedAt.String()),
	})

	c.JSON(201, event)
}
