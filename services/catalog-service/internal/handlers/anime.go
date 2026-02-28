package handlers

import (
	"catalog-service/internal/models"
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AnimeHandler struct {
	db     *gorm.DB
	reader *kafka.Reader
}

func (h *AnimeHandler) SetSources(db* gorm.DB, reader *kafka.Reader){
	h.db = db
	h.reader = reader
}

func (h *AnimeHandler) GetAnimes(c *gin.Context) {
	var list []models.Anime
	h.db.Find(&list)
	c.JSON(200, list)

}

func (h *AnimeHandler) GetAnime(c *gin.Context) {
	idParam := c.Param("id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid uuid"})
		return
	}

	var anime models.Anime
	if h.db.First(&anime, id).Error != nil {
		c.JSON(404, gin.H{"error": "not found"})
		return
	}
	c.JSON(200, anime)
}

func (h *AnimeHandler) HandleKafkaEvents() {
	for {
		msg, err := h.reader.ReadMessage(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		animeId, err := uuid.ParseBytes(msg.Key)
		if err != nil {
			log.Printf("invalid UUID key: %v", err)
			continue
		}

		updatedAt, err := time.Parse(time.RFC3339, string(msg.Value))
		if err != nil {
			log.Printf("invalid timestamp value: %v", err)
			continue
		}

		err = h.db.Model(&models.Anime{}).
			Where("id = ?", animeId).
			Updates(map[string]any{
				"watch_count": gorm.Expr("watch_count + 1"),
				"updated_at":  updatedAt,
			}).Error

		if err != nil {
			log.Printf("failed to update anime %s: %v", animeId, err)
		}
	}
}
