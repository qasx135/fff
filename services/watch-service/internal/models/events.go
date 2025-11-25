package models

import "gorm.io/gorm"

type WatchEvent struct {
	gorm.Model
	UserID  string `json:"user_id"`
	AnimeID uint   `json:"anime_id"`
	Episode int    `json:"episode"`
}