package database

import (
	"fmt"
	"github.com/helioLJ/tweetvault/config"
	"github.com/helioLJ/tweetvault/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto-migrate the schema
	err = db.AutoMigrate(
		&models.Bookmark{},
		&models.Media{},
		&models.Tag{},
		&models.BookmarkTag{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}