package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/helioLJ/tweetvault/internal/models"
	"github.com/helioLJ/tweetvault/internal/services"
	"gorm.io/gorm"
)

type BookmarkHandler struct {
	service *services.BookmarkService
	db      *gorm.DB
}

func NewBookmarkHandler(db *gorm.DB) *BookmarkHandler {
	return &BookmarkHandler{
		service: services.NewBookmarkService(db),
		db:      db,
	}
}

// List returns all bookmarks with optional filtering
func (h *BookmarkHandler) List(c *gin.Context) {
	tag := c.Query("tag")
	search := c.Query("search")
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "12")
	showArchived := c.DefaultQuery("archived", "false") == "true"

	bookmarks, total, err := h.service.List(tag, search, page, limit, showArchived)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"bookmarks": bookmarks,
		"total":     total,
	})
}

// Get returns a single bookmark by ID
func (h *BookmarkHandler) Get(c *gin.Context) {
	bookmark, err := h.service.Get(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bookmark not found"})
		return
	}

	c.JSON(http.StatusOK, bookmark)
}

// Update updates a bookmark's tags
func (h *BookmarkHandler) Update(c *gin.Context) {
	var input struct {
		Tags []string `json:"tags"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateTags(c.Param("id"), input.Tags); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tags updated successfully"})
}

// Delete removes a bookmark
func (h *BookmarkHandler) Delete(c *gin.Context) {
	if err := h.service.Delete(c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bookmark deleted successfully"})
}

// GetStatistics returns overall statistics about bookmarks and tags
func (h *BookmarkHandler) GetStatistics(c *gin.Context) {
	var stats struct {
		TotalBookmarks    int64 `json:"total_bookmarks"`
		ActiveBookmarks   int64 `json:"active_bookmarks"`
		ArchivedBookmarks int64 `json:"archived_bookmarks"`
		TotalTags        int64 `json:"total_tags"`
		TopTags          []struct {
			Name  string `json:"name"`
			Count int64  `json:"count"`
		} `json:"top_tags"`
	}

	// Initialize TopTags as empty slice instead of nil
	stats.TopTags = []struct {
		Name  string `json:"name"`
		Count int64  `json:"count"`
	}{}

	// Get total bookmarks
	h.db.Model(&models.Bookmark{}).Count(&stats.TotalBookmarks)

	// Get active bookmarks
	h.db.Model(&models.Bookmark{}).Where("archived = ?", false).Count(&stats.ActiveBookmarks)

	// Get archived bookmarks
	h.db.Model(&models.Bookmark{}).Where("archived = ?", true).Count(&stats.ArchivedBookmarks)

	// Get total tags
	h.db.Model(&models.Tag{}).Count(&stats.TotalTags)

	// Get top 5 tags with their counts
	h.db.Raw(`
		SELECT t.name, COUNT(bt.bookmark_id) as count 
		FROM tags t 
		JOIN bookmark_tags bt ON t.id = bt.tag_id 
		JOIN bookmarks b ON bt.bookmark_id = b.id
		WHERE b.archived = false
		GROUP BY t.id, t.name 
		ORDER BY count DESC 
		LIMIT 5
	`).Scan(&stats.TopTags)

	c.JSON(http.StatusOK, stats)
}

func (h *BookmarkHandler) ToggleArchive(c *gin.Context) {
	if err := h.service.ToggleArchive(c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Bookmark archive status toggled successfully"})
}
