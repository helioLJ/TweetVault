package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"

	"archive/zip"

	"github.com/gin-gonic/gin"
	"github.com/helioLJ/tweetvault/internal/models"
	"gorm.io/gorm"
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
		fmt.Printf("Error getting JSON file: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "No JSON file provided"})
		return
	}

	// Get the ZIP file
	zipFile, err := c.FormFile("zipFile")
	if err != nil {
		fmt.Printf("Error getting ZIP file: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "No ZIP file provided"})
		return
	}

	// Process JSON file
	bookmarks, err := h.processJSONFile(jsonFile)
	if err != nil {
		fmt.Printf("Error processing JSON file: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Begin transaction
	tx := h.db.Begin()

	// Process each bookmark
	for i, bookmark := range bookmarks {
		// Create or update bookmark
		if err := h.processBookmark(tx, bookmark, zipFile); err != nil {
			fmt.Printf("Error processing bookmark %d: %v\n", i, err)
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		fmt.Printf("Error committing transaction: %v\n", err)
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
			// Log the error but continue processing
			fmt.Printf("Warning: Could not extract media file %s: %v\n", mediaFileName, err)
			continue
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
	// Try parsing with timezone offset format (with space before timezone)
	t, err := time.Parse("2006-01-02 15:04:05 -0700", timeStr)
	if err == nil {
		return t
	}

	// Try parsing with timezone offset format (without space)
	t, err = time.Parse("2006-01-02 15:04:05-0700", timeStr)
	if err == nil {
		return t
	}

	// Try parsing with timezone offset format (with space and colon in timezone)
	t, err = time.Parse("2006-01-02 15:04:05 -07:00", timeStr)
	if err == nil {
		return t
	}

	// Fallback to parsing without timezone
	t, err = time.Parse("2006-01-02 15:04:05", timeStr)
	if err == nil {
		return t
	}

	// If all parsing attempts fail, log the error and return current time
	fmt.Printf("Error parsing time '%s': %v\n", timeStr, err)
	return time.Now()
}

func generateMediaFileName(screenName, tweetID, mediaType string, index int) string {
	// Extract creation date from tweet ID to match the file naming convention
	// Twitter IDs contain a timestamp that we can use
	tweetIDInt, _ := strconv.ParseInt(tweetID, 10, 64)
	timestamp := time.Unix((tweetIDInt>>22)/1000+1288834974657/1000, 0)
	dateStr := timestamp.Format("20060102") // Format as YYYYMMDD

	return fmt.Sprintf("%s_%s_%s_%d_%s%s",
		screenName,
		tweetID,
		mediaType,
		index,
		dateStr,
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
	// Open the uploaded zip file
	src, err := zipFile.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open zip file: %w", err)
	}
	defer src.Close()

	// Create a temporary file to store the zip content
	tempFile, err := os.CreateTemp("", "upload-*.zip")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Copy zip content to temp file
	if _, err := io.Copy(tempFile, src); err != nil {
		return nil, fmt.Errorf("failed to copy zip content: %w", err)
	}

	// Open the zip archive
	reader, err := zip.OpenReader(tempFile.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to open zip archive: %w", err)
	}
	defer reader.Close()

	// Debug: Print all files in the ZIP
	fmt.Println("Files in ZIP:")
	for _, f := range reader.File {
		fmt.Printf("- %s\n", f.Name)
	}
	fmt.Printf("Looking for file: %s\n", fileName)

	// Find and extract the requested file
	for _, file := range reader.File {
		if file.Name == fileName {
			rc, err := file.Open()
			if err != nil {
				return nil, fmt.Errorf("failed to open file in zip: %w", err)
			}
			defer rc.Close()

			return io.ReadAll(rc)
		}
	}

	return nil, fmt.Errorf("file %s not found in zip archive", fileName)
}
