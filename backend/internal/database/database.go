package database

import (
	"fmt"
	"log"

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

	// Create materialized view
	if err := CreateBookmarkView(db); err != nil {
		log.Printf("Warning: Failed to create materialized view: %v", err)
	}

	// Ensure standard tags exist
	if err := ensureStandardTags(db); err != nil {
		return nil, err
	}

	// Add join table configuration
	db.SetupJoinTable(&models.Bookmark{}, "Tags", &models.BookmarkTag{})
	db.SetupJoinTable(&models.Tag{}, "Bookmarks", &models.BookmarkTag{})

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

// Move CreateBookmarkView here from migrations/create_bookmark_view.go
func CreateBookmarkView(db *gorm.DB) error {
	// Drop existing view if it exists
	db.Exec(`DROP MATERIALIZED VIEW IF EXISTS bookmark_views`)

	// Create the materialized view
	return db.Exec(`
		CREATE MATERIALIZED VIEW bookmark_views AS
		SELECT 
			b.id,
			b.created_at,
			b.full_text,
			b.screen_name,
			b.name,
			b.profile_image_url,
			b.favorite_count,
			b.retweet_count,
			b.views_count,
			b.url,
			b.archived,
			COALESCE(
				(
					SELECT json_agg(json_build_object(
						'id', m.id,
						'type', m.type,
						'url', m.url,
						'thumbnail', m.thumbnail,
						'original', m.original
					))
					FROM media m
					WHERE m.tweet_id = b.id
				),
				'[]'::json
			) as media_json,
			COALESCE(
				(
					SELECT json_agg(json_build_object(
						'id', t.id,
						'name', t.name,
						'completed', bt.completed
					))
					FROM tags t
					JOIN bookmark_tags bt ON bt.tag_id = t.id
					WHERE bt.bookmark_id = b.id
				),
				'[]'::json
			) as tags_json
		FROM bookmarks b
		GROUP BY b.id
		ORDER BY b.created_at DESC;

		CREATE UNIQUE INDEX idx_bookmark_views_id ON bookmark_views(id);
		CREATE INDEX idx_bookmark_views_archived_created ON bookmark_views(archived, created_at DESC);
	`).Error
}
