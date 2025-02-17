package models

import (
	"time"
)

type Media struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	TweetID   string    `gorm:"type:varchar(30);index:idx_media_tweet" json:"tweet_id"`
	Type      string    `gorm:"type:varchar(20)" json:"type"`       // video, photo
	URL       string    `gorm:"type:text" json:"url"`               // Original Twitter URL
	Thumbnail string    `gorm:"type:text" json:"thumbnail"`         // Twitter thumbnail URL
	Original  string    `gorm:"type:text" json:"original"`          // Twitter original media URL
	FileData  []byte    `gorm:"type:bytea" json:"-"`                // Actual media file content
	FileName  string    `gorm:"type:varchar(255)" json:"file_name"` // Original filename from data/media
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName specifies the table name for the Media model
func (Media) TableName() string {
	return "media"
}
