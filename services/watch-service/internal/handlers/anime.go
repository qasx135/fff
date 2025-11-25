package handlers

import (
	"watch-service/internal/models"

	"github.com/gin-gonic/gin"
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

	if event.UserID == "" || event.AnimeID == 0 {
		c.JSON(400, gin.H{"error": "user_id and anime_id required"})
		return
	}

	h.db.Create(&event)
	c.JSON(201, event)
}
