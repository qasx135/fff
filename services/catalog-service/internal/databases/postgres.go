package databases

import (
	"catalog-service/internal/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

func InitDB(dsn string) {
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Не удалось подключиться к БД:", err)
	}
	DB.AutoMigrate(&models.Anime{})
	DB.FirstOrCreate(&models.Anime{}, models.Anime{
		Title:      "Cowboy Bebop",
		Genre:      "Sci-Fi",
		Episodes:   26,
		WatchCount: 0,
	})
}

func HealthCheck() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}
