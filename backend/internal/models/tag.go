package models

import (
	"time"
)

type Tag struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(50);unique" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	Bookmarks []Bookmark `gorm:"many2many:bookmark_tags" json:"bookmarks,omitempty"`
}

type BookmarkTag struct {
	BookmarkID string    `gorm:"primaryKey;type:varchar(30)" json:"bookmark_id"`
	TagID      uint      `gorm:"primaryKey" json:"tag_id"`
	CreatedAt  time.Time `json:"created_at"`
}