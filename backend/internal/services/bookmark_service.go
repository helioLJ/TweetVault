package services

import (
	"github.com/helioLJ/tweetvault/internal/models"
	"gorm.io/gorm"
	"strconv"
)

type BookmarkService struct {
	db *gorm.DB
}

func NewBookmarkService(db *gorm.DB) *BookmarkService {
	return &BookmarkService{db: db}
}

func (s *BookmarkService) List(tag string, search string, page string, limit string) ([]models.Bookmark, int64, error) {
	var bookmarks []models.Bookmark
	var total int64
	query := s.db.Model(&models.Bookmark{})

	if tag != "" {
		query = query.Joins("JOIN bookmark_tags ON bookmarks.id = bookmark_tags.bookmark_id").
			Joins("JOIN tags ON bookmark_tags.tag_id = tags.id").
			Where("tags.name = ?", tag)
	}

	if search != "" {
		query = query.Where("full_text ILIKE ? OR name ILIKE ? OR screen_name ILIKE ?", 
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	query.Count(&total)

	// Parse pagination params
	pageNum, _ := strconv.Atoi(page)
	limitNum, _ := strconv.Atoi(limit)
	offset := (pageNum - 1) * limitNum

	err := query.Preload("Tags").Preload("Media").
		Offset(offset).Limit(limitNum).
		Find(&bookmarks).Error

	return bookmarks, total, err
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