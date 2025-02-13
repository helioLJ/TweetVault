package models

import (
	"time"
)

type List struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	Name      string     `gorm:"type:varchar(100);not null" json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	Bookmarks []Bookmark `gorm:"many2many:list_bookmarks;constraint:OnDelete:CASCADE" json:"bookmarks,omitempty"`
}

type ListBookmark struct {
	ListID     uint      `gorm:"primaryKey" json:"list_id"`
	BookmarkID string    `gorm:"primaryKey;type:varchar(30)" json:"bookmark_id"`
	Position   int       `gorm:"not null" json:"position"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}