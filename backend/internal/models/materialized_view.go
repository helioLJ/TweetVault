package models

import "time"

// BookmarkView represents the materialized view structure
type BookmarkView struct {
	ID              string          `gorm:"primaryKey;column:id" json:"id"`
	CreatedAt       time.Time       `json:"created_at"`
	FullText        string          `json:"full_text"`
	ScreenName      string          `json:"screen_name"`
	Name            string          `json:"name"`
	ProfileImageURL string          `json:"profile_image_url"`
	FavoriteCount   int             `json:"favorite_count"`
	RetweetCount    int             `json:"retweet_count"`
	ViewsCount      int             `json:"views_count"`
	URL             string          `json:"url"`
	Archived        bool            `json:"archived"`
	Media           []Media         `gorm:"-" json:"media"`    // Will be populated from JSON
	Tags            []TagWithStatus `gorm:"-" json:"tags"`     // Will be populated from JSON
	MediaJSON       string          `gorm:"column:media_json"` // Stored as JSON string
	TagsJSON        string          `gorm:"column:tags_json"`  // Stored as JSON string
}

type TagWithStatus struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Completed bool   `json:"completed"`
}

// TableName specifies the materialized view name
func (BookmarkView) TableName() string {
	return "bookmark_views"
}
