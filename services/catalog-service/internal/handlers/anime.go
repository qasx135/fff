package handlers

import (
	"catalog-service/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AnimeHandler struct {
	db *gorm.DB
}

func (h *AnimeHandler) GetAnime(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		var list []models.Anime
		h.db.Find(&list)
		c.JSON(200, list)
		return
	}

	var anime models.Anime
	if h.db.First(&anime, id).Error != nil {
		c.JSON(404, gin.H{"error": "not found"})
		return
	}
	c.JSON(200, anime)
}
