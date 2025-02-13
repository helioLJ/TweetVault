package services

import (
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
	err := query.Preload("Media").
		Preload("Tags").
		Order("created_at DESC").
		Offset(offset).
		Limit(limitNum).
		Find(&bookmarks).Error

	if err != nil {
		log.Printf("Error executing final query: %v", err)
		return nil, 0, err
	}

	// After loading, set the completed status for each tag
	for i := range bookmarks {
		for j := range bookmarks[i].Tags {
			var bookmarkTag models.BookmarkTag
			if err := s.db.Where("bookmark_id = ? AND tag_id = ?",
				bookmarks[i].ID, bookmarks[i].Tags[j].ID).
				First(&bookmarkTag).Error; err == nil {
				bookmarks[i].Tags[j].Completed = bookmarkTag.Completed
			}
		}
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

	if err := tx.Exec("DELETE FROM bookmark_tags WHERE bookmark_id = ?", id).Error; err != nil {
		tx.Rollback()
		return err
	}

	for _, tagName := range tags {
		var tag models.Tag
		if err := tx.FirstOrCreate(&tag, models.Tag{Name: tagName}).Error; err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.Create(&models.BookmarkTag{
			BookmarkID: id,
			TagID:      tag.ID,
		}).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (s *BookmarkService) Delete(id string) error {
	return s.db.Delete(&models.Bookmark{}, "id = ?", id).Error
}

func (s *BookmarkService) ToggleArchive(id string) error {
	return s.db.Model(&models.Bookmark{}).Where("id = ?", id).
		Update("archived", gorm.Expr("NOT archived")).Error
}
