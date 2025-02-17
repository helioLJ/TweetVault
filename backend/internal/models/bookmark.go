package models

import (
	"database/sql"
	"encoding/json"
	"time"
)

type Bookmark struct {
	ID              string          `gorm:"primaryKey;column:id;type:varchar(30)" json:"id"`
	CreatedAt       time.Time       `json:"created_at" gorm:"index:idx_archived_createdAt"`
	FullText        string          `gorm:"type:text" json:"full_text"`
	ScreenName      string          `gorm:"type:varchar(50)" json:"screen_name"`
	Name            string          `gorm:"type:varchar(100)" json:"name"`
	ProfileImageURL string          `gorm:"type:text" json:"profile_image_url"`
	InReplyTo       sql.NullString  `gorm:"type:varchar(30)" json:"in_reply_to"`
	RetweetedStatus sql.NullString  `gorm:"type:varchar(30)" json:"retweeted_status"`
	QuotedStatus    sql.NullString  `gorm:"type:varchar(30)" json:"quoted_status"`
	FavoriteCount   int             `json:"favorite_count"`
	RetweetCount    int             `json:"retweet_count"`
	BookmarkCount   int             `json:"bookmark_count"`
	QuoteCount      int             `json:"quote_count"`
	ReplyCount      int             `json:"reply_count"`
	ViewsCount      int             `json:"views_count"`
	Favorited       bool            `json:"favorited"`
	Retweeted       bool            `json:"retweeted"`
	Bookmarked      bool            `json:"bookmarked"`
	URL             string          `gorm:"type:text" json:"url"`
	Metadata        json.RawMessage `gorm:"type:jsonb" json:"metadata"`
	Media           []Media         `gorm:"foreignKey:TweetID" json:"media"`
	Tags            []Tag           `gorm:"many2many:bookmark_tags" json:"tags"`
	UpdatedAt       time.Time       `gorm:"autoUpdateTime" json:"-"`
	Archived        bool            `json:"archived" gorm:"default:false;index:idx_archived_createdAt"`
}

// TableName specifies the table name for the Bookmark model
func (Bookmark) TableName() string {
	return "bookmarks"
}
