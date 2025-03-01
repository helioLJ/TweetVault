package models

import (
	"time"
)

type Tag struct {
	ID           uint          `gorm:"primaryKey" json:"id"`
	Name         string        `gorm:"type:varchar(50);unique" json:"name"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
	Completed    bool          `json:"completed"`
	Bookmarks    []Bookmark    `gorm:"many2many:bookmark_tags" json:"bookmarks,omitempty"`
	BookmarkTags []BookmarkTag `gorm:"foreignKey:TagID" json:"-"`
}

type BookmarkTag struct {
	BookmarkID string    `gorm:"primaryKey;type:varchar(30);index:idx_bookmark_id" json:"bookmark_id"`
	TagID      uint      `gorm:"primaryKey;index:idx_tag_id" json:"tag_id"`
	CreatedAt  time.Time `json:"created_at"`
	Completed  bool      `gorm:"default:false;index:idx_completed" json:"completed"`
}
