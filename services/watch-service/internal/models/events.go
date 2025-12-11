package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WatchEvent struct {
	gorm.Model
	UserID  string    `json:"user_id"`
	AnimeID uuid.UUID `gorm:"type:uuid" json:"anime_id"`
	Episode int       `json:"episode"`
}
