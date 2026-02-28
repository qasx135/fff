package handlers

import (
	"errors"
	"time"
	internal_kafka "watch-service/internal/kafka"
	"watch-service/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
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
	ctx := c.Request.Context()
	tracer := otel.Tracer("watch-handler")

	ctx, span := tracer.Start(ctx, "WatchHandler")
	defer span.End()

	var dto WatchEventDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	span.SetAttributes(
		attribute.String("user.id", dto.UserID),
		attribute.String("anime.id", dto.AnimeID.String()),
		attribute.Int("episode", dto.Episode),
	)

	if dto.UserID == "" || dto.AnimeID == uuid.Nil {
		err := errors.New("user_id and anime_id required")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	dbCtx, dbSpan := tracer.Start(ctx, "db.insert_watch_event")
	event := models.WatchEvent{
		UserID:  dto.UserID,
		AnimeID: dto.AnimeID,
		Episode: dto.Episode,
	}

	if err := h.db.WithContext(dbCtx).Create(&event).Error; err != nil {
		dbSpan.RecordError(err)
		dbSpan.SetStatus(codes.Error, "db error")
		dbSpan.End()

		c.JSON(500, gin.H{"error": "db error"})
		return
	}
	dbSpan.End()

	kafkaCtx, kafkaSpan := tracer.Start(ctx, "kafka.produce_watch_event")

	err := internal_kafka.WriteMessage(
		kafkaCtx,
		kafka.Message{
			Key:   []byte(event.AnimeID.String()),
			Value: []byte(event.UpdatedAt.Format(time.RFC3339)),
		},
	)

	if err != nil {
		kafkaSpan.RecordError(err)
		kafkaSpan.SetStatus(codes.Error, "kafka produce failed")
	}
	kafkaSpan.End()

	c.JSON(201, event)
}
