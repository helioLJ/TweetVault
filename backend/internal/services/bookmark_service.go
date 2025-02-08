package services

import (
	"github.com/helioLJ/tweetvault/internal/models"
	"gorm.io/gorm"
)

type BookmarkService struct {
	db *gorm.DB
}

func NewBookmarkService(db *gorm.DB) *BookmarkService {
	return &BookmarkService{db: db}
}

func (s *BookmarkService) List(tag string) ([]models.Bookmark, int64, error) {
	var bookmarks []models.Bookmark
	query := s.db.Preload("Media").Preload("Tags")

	if tag != "" {
		query = query.Joins("JOIN bookmark_tags ON bookmarks.id = bookmark_tags.bookmark_id").
			Joins("JOIN tags ON bookmark_tags.tag_id = tags.id").
			Where("tags.name = ?", tag)
	}

	var total int64
	if err := query.Model(&models.Bookmark{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Find(&bookmarks).Error; err != nil {
		return nil, 0, err
	}

	return bookmarks, total, nil
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