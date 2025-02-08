package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"github.com/helioLJ/tweetvault/internal/models"
)

type UploadHandler struct {
	db *gorm.DB
}

func NewUploadHandler(db *gorm.DB) *UploadHandler {
	return &UploadHandler{db: db}
}

type TwitterBookmark struct {
	ID              string          `json:"id"`
	CreatedAt       string          `json:"created_at"`
	FullText        string          `json:"full_text"`
	ScreenName      string          `json:"screen_name"`
	Name            string          `json:"name"`
	ProfileImageURL string          `json:"profile_image_url"`
	InReplyTo       *string         `json:"in_reply_to"`
	RetweetedStatus *string         `json:"retweeted_status"`
	QuotedStatus    *string         `json:"quoted_status"`
	FavoriteCount   int             `json:"favorite_count"`
	RetweetCount    int             `json:"retweet_count"`
	BookmarkCount   int             `json:"bookmark_count"`
	QuoteCount      int             `json:"quote_count"`
	ReplyCount      int             `json:"reply_count"`
	ViewsCount      int             `json:"views_count"`
	Favorited       bool            `json:"favorited"`
	Retweeted       bool            `json:"retweeted"`
	Bookmarked      bool            `json:"bookmarked"`
	URL             string          `json:"url"`
	Metadata        json.RawMessage `json:"metadata"`
	Media           []struct {
		Type      string `json:"type"`
		URL       string `json:"url"`
		Thumbnail string `json:"thumbnail"`
		Original  string `json:"original"`
	} `json:"media"`
}

func (h *UploadHandler) HandleUpload(c *gin.Context) {
	// Get the JSON file
	jsonFile, err := c.FormFile("jsonFile")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No JSON file provided"})
		return
	}

	// Get the ZIP file
	zipFile, err := c.FormFile("zipFile")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No ZIP file provided"})
		return
	}

	// Process JSON file
	bookmarks, err := h.processJSONFile(jsonFile)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Begin transaction
	tx := h.db.Begin()

	// Process each bookmark
	for _, bookmark := range bookmarks {
		// Create or update bookmark
		if err := h.processBookmark(tx, bookmark, zipFile); err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Upload processed successfully",
		"count":   len(bookmarks),
	})
}

func (h *UploadHandler) processJSONFile(file *multipart.FileHeader) ([]TwitterBookmark, error) {
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	bytes, err := io.ReadAll(src)
	if err != nil {
		return nil, err
	}

	var bookmarks []TwitterBookmark
	if err := json.Unmarshal(bytes, &bookmarks); err != nil {
		return nil, err
	}

	return bookmarks, nil
}

func (h *UploadHandler) processBookmark(tx *gorm.DB, tb TwitterBookmark, zipFile *multipart.FileHeader) error {
	bookmark := models.Bookmark{
		ID:              tb.ID,
		CreatedAt:       parseTwitterTime(tb.CreatedAt),
		FullText:        tb.FullText,
		ScreenName:      tb.ScreenName,
		Name:            tb.Name,
		ProfileImageURL: tb.ProfileImageURL,
		FavoriteCount:   tb.FavoriteCount,
		RetweetCount:    tb.RetweetCount,
		BookmarkCount:   tb.BookmarkCount,
		QuoteCount:      tb.QuoteCount,
		ReplyCount:      tb.ReplyCount,
		ViewsCount:      tb.ViewsCount,
		Favorited:       tb.Favorited,
		Retweeted:       tb.Retweeted,
		Bookmarked:      tb.Bookmarked,
		URL:             tb.URL,
		Metadata:        tb.Metadata,
	}

	// Create or update bookmark
	if err := tx.Save(&bookmark).Error; err != nil {
		return err
	}

	// Process media files from ZIP
	for i, m := range tb.Media {
		mediaFileName := generateMediaFileName(tb.ScreenName, tb.ID, m.Type, i+1)
		mediaData, err := extractFileFromZip(zipFile, mediaFileName)
		if err != nil {
			return err
		}

		media := models.Media{
			TweetID:   tb.ID,
			Type:      m.Type,
			URL:       m.URL,
			Thumbnail: m.Thumbnail,
			Original:  m.Original,
			FileData:  mediaData,
			FileName:  mediaFileName,
		}

		if err := tx.Create(&media).Error; err != nil {
			return err
		}
	}

	return nil
}

// Helper functions (to be implemented)
func parseTwitterTime(timeStr string) time.Time {
	// Implement time parsing logic
	t, _ := time.Parse("2006-01-02 15:04:05 -0700", timeStr)
	return t
}

func generateMediaFileName(screenName, tweetID, mediaType string, index int) string {
	return fmt.Sprintf("%s_%s_%s_%d_%s%s",
		screenName,
		tweetID,
		mediaType,
		index,
		time.Now().Format("20060102"),
		getExtensionForMediaType(mediaType))
}

func getExtensionForMediaType(mediaType string) string {
	switch mediaType {
	case "video":
		return ".mp4"
	case "photo":
		return ".jpg"
	default:
		return ""
	}
}

func extractFileFromZip(zipFile *multipart.FileHeader, fileName string) ([]byte, error) {
	// Implement ZIP extraction logic
	// This will need to use archive/zip package to extract the specific file
	// and return its contents as []byte
	return nil, nil
}
