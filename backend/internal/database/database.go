package database

import (
	"fmt"

	"github.com/helioLJ/tweetvault/config"
	"github.com/helioLJ/tweetvault/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Fixed tags that should always exist in the system
var standardTags = []string{"To do", "To read"}

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

	// Ensure standard tags exist
	if err := ensureStandardTags(db); err != nil {
		return nil, err
	}

	return db, nil
}

// ensureStandardTags creates the standard tags if they don't exist
func ensureStandardTags(db *gorm.DB) error {
	for _, tagName := range standardTags {
		var tag models.Tag
		result := db.Where("name = ?", tagName).First(&tag)
		if result.Error == gorm.ErrRecordNotFound {
			// Create the tag if it doesn't exist
			tag = models.Tag{Name: tagName}
			if err := db.Create(&tag).Error; err != nil {
				return fmt.Errorf("failed to create standard tag %s: %w", tagName, err)
			}
		} else if result.Error != nil {
			return fmt.Errorf("error checking for standard tag %s: %w", tagName, result.Error)
		}
	}
	return nil
}
