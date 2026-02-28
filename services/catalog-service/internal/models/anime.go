package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Anime struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey"`
	CreatedAt   time.Time  `gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt   *time.Time `gorm:"type:timestamptz"`
	Title       string
	Description string
	Genre       string
	Episodes    int
	WatchCount  int
}


func (p *Anime) BeforeCreate(tx *gorm.DB) (err error) {
    if p.ID == uuid.Nil {
        p.ID = uuid.New()
    }
    return
}