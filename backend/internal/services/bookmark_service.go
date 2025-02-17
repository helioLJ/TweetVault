package services

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/helioLJ/tweetvault/internal/models"
	"gorm.io/gorm"
)

type BookmarkService struct {
	db *gorm.DB
}

func NewBookmarkService(db *gorm.DB) *BookmarkService {
	return &BookmarkService{db: db}
}

func (s *BookmarkService) List(tag string, search string, page string, limit string, showArchived bool) ([]models.Bookmark, int64, error) {
	var bookmarks []models.Bookmark
	var total int64

	// Start with a new query
	query := s.db.Model(&models.Bookmark{})

	// Debug logging
	log.Printf("List Bookmarks - Input parameters:")
	log.Printf("  showArchived: %v", showArchived)
	log.Printf("  tag: %v", tag)
	log.Printf("  search: %v", search)
	log.Printf("  page: %v", page)
	log.Printf("  limit: %v", limit)

	// Apply archived filter
	query = query.Where("archived = ?", showArchived)

	// Apply tag filter if provided
	if tag != "" {
		query = query.Joins("LEFT JOIN bookmark_tags ON bookmarks.id = bookmark_tags.bookmark_id").
			Joins("LEFT JOIN tags ON bookmark_tags.tag_id = tags.id").
			Where("tags.name = ?", tag)
	}

	// Apply search filter if provided
	if search != "" {
		query = query.Where("full_text ILIKE ? OR name ILIKE ? OR screen_name ILIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// Get total count before pagination
	if err := query.Count(&total).Error; err != nil {
		log.Printf("Error counting total records: %v", err)
		return nil, 0, err
	}

	// Parse pagination params
	pageNum, _ := strconv.Atoi(page)
	limitNum, _ := strconv.Atoi(limit)
	offset := (pageNum - 1) * limitNum

	// Execute final query with preloads and pagination
	err := query.
		Preload("Media").
		Preload("Tags", func(db *gorm.DB) *gorm.DB {
			return db.Select("tags.id, tags.name, bookmark_tags.completed").
				Joins("LEFT JOIN bookmark_tags ON bookmark_tags.tag_id = tags.id")
		}).
		Order("created_at DESC").
		Offset(offset).
		Limit(limitNum).
		Find(&bookmarks).Error

	if err != nil {
		log.Printf("Error executing final query: %v", err)
		return nil, 0, err
	}

	// Debug output
	log.Printf("Total matching records: %d", total)
	log.Printf("Found %d bookmarks in current page", len(bookmarks))
	if len(bookmarks) > 0 {
		for i := 0; i < min(3, len(bookmarks)); i++ {
			log.Printf("Sample bookmark %d - ID: %s, Archived: %v",
				i+1, bookmarks[i].ID, bookmarks[i].Archived)
			// Add debug logging for tags and their completion status
			for _, tag := range bookmarks[i].Tags {
				log.Printf("  Tag: %s, Completed: %v", tag.Name, tag.Completed)
			}
		}
	}

	return bookmarks, total, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (s *BookmarkService) Get(id string) (*models.Bookmark, error) {
	var bookmark models.Bookmark
	if err := s.db.Preload("Media").Preload("Tags").First(&bookmark, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &bookmark, nil
}

func (s *BookmarkService) UpdateTags(id string, tags []string) error {
	tx := s.db.Begin()

	// Instead of deleting all tags, get existing ones first
	var existingBookmarkTags []models.BookmarkTag
	if err := tx.Where("bookmark_id = ?", id).Find(&existingBookmarkTags).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Create a map of existing tag completion status
	existingCompletionStatus := make(map[string]bool)
	for _, bt := range existingBookmarkTags {
		var tag models.Tag
		if err := tx.First(&tag, bt.TagID).Error; err != nil {
			continue
		}
		existingCompletionStatus[tag.Name] = bt.Completed
	}

	// Now delete existing tags
	if err := tx.Where("bookmark_id = ?", id).Delete(&models.BookmarkTag{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Add all tags, preserving completion status for existing ones
	for _, tagName := range tags {
		var tag models.Tag
		if err := tx.FirstOrCreate(&tag, models.Tag{Name: tagName}).Error; err != nil {
			tx.Rollback()
			return err
		}

		completed := existingCompletionStatus[tagName] // Get existing completion status if any
		if err := tx.Create(&models.BookmarkTag{
			BookmarkID: id,
			TagID:      tag.ID,
			Completed:  completed,
		}).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (s *BookmarkService) Delete(id string) error {
	// Start a transaction
	return s.db.Transaction(func(tx *gorm.DB) error {
		// First delete associated records in bookmark_tags
		if err := tx.Where("bookmark_id = ?", id).Delete(&models.BookmarkTag{}).Error; err != nil {
			return err
		}

		// Then delete associated media
		if err := tx.Where("tweet_id = ?", id).Delete(&models.Media{}).Error; err != nil {
			return err
		}

		// Finally delete the bookmark
		if err := tx.Delete(&models.Bookmark{}, "id = ?", id).Error; err != nil {
			return err
		}

		return nil
	})
}

func (s *BookmarkService) ToggleArchive(id string) error {
	return s.db.Model(&models.Bookmark{}).Where("id = ?", id).
		Update("archived", gorm.Expr("NOT archived")).Error
}

// RefreshView refreshes the materialized view
func (s *BookmarkService) RefreshView() error {
	return s.db.Exec("REFRESH MATERIALIZED VIEW CONCURRENTLY bookmark_views").Error
}

// ListFromView gets bookmarks from the materialized view
func (s *BookmarkService) ListFromView(tag string, search string, page string, limit string, showArchived bool) ([]models.BookmarkView, int64, error) {
	var bookmarks []models.BookmarkView
	var total int64

	query := s.db.Model(&models.BookmarkView{})

	// Apply archived filter
	query = query.Where("archived = ?", showArchived)

	// Apply tag filter if provided
	if tag != "" {
		query = query.Where("tags_json::jsonb @> ?", fmt.Sprintf(`[{"name":"%s"}]`, tag))
	}

	// Apply search filter if provided
	if search != "" {
		query = query.Where("full_text ILIKE ? OR name ILIKE ? OR screen_name ILIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Parse pagination params
	pageNum, _ := strconv.Atoi(page)
	limitNum, _ := strconv.Atoi(limit)
	offset := (pageNum - 1) * limitNum

	// Execute final query
	err := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(limitNum).
		Find(&bookmarks).Error

	if err != nil {
		return nil, 0, err
	}

	// Parse JSON fields into structs
	for i := range bookmarks {
		if bookmarks[i].MediaJSON != "" {
			if err := json.Unmarshal([]byte(bookmarks[i].MediaJSON), &bookmarks[i].Media); err != nil {
				log.Printf("Error unmarshaling media JSON: %v", err)
			}
		}
		if bookmarks[i].TagsJSON != "" {
			if err := json.Unmarshal([]byte(bookmarks[i].TagsJSON), &bookmarks[i].Tags); err != nil {
				log.Printf("Error unmarshaling tags JSON: %v", err)
			}
		}
	}

	return bookmarks, total, nil
}
