package models

import "gorm.io/gorm"

type Anime struct {
	gorm.Model
	Title       string
	Description string
	Genre       string
	Episodes    int
	WatchCount  int
}
