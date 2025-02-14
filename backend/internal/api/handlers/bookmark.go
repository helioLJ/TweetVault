package handlers

import (
	"net/http"

	"log"

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

	// Add debug logging
	archivedParam := c.Query("archived")
	log.Printf("Received archived parameter: %v", archivedParam)

	// Parse boolean more explicitly
	showArchived := false
	if archivedParam == "true" {
		showArchived = true
	}

	log.Printf("Parsed showArchived value: %v", showArchived)

	bookmarks, total, err := h.service.List(tag, search, page, limit, showArchived)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Load tags with completion status for each bookmark
	for i := range bookmarks {
		if err := h.db.Preload("Media").
			Preload("Tags", func(db *gorm.DB) *gorm.DB {
				return db.Select("tags.*, bookmark_tags.completed").
					Joins("LEFT JOIN bookmark_tags ON bookmark_tags.tag_id = tags.id AND bookmark_tags.bookmark_id = ?", bookmarks[i].ID)
			}).
			First(&bookmarks[i], bookmarks[i].ID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"bookmarks": bookmarks,
		"total":     total,
	})
}

// Get returns a single bookmark by ID
func (h *BookmarkHandler) Get(c *gin.Context) {
	var bookmark models.Bookmark
	if err := h.db.Preload("Media").
		Preload("Tags", func(db *gorm.DB) *gorm.DB {
			return db.Select("tags.*, bookmark_tags.completed").
				Joins("LEFT JOIN bookmark_tags ON bookmark_tags.tag_id = tags.id AND bookmark_tags.bookmark_id = ?", c.Param("id"))
		}).
		First(&bookmark, c.Param("id")).Error; err != nil {
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

// In the statistics handler
type TagStats struct {
	Name           string `json:"name"`
	Count          int64  `json:"count"`
	CompletedCount int64  `json:"completed_count"`
}

func (h *BookmarkHandler) GetStatistics(c *gin.Context) {
	var stats struct {
		TotalBookmarks    int64      `json:"total_bookmarks"`
		ActiveBookmarks   int64      `json:"active_bookmarks"`
		ArchivedBookmarks int64      `json:"archived_bookmarks"`
		TotalTags         int64      `json:"total_tags"`
		TopTags           []TagStats `json:"top_tags"`
	}

	// Initialize TopTags as empty slice instead of nil
	stats.TopTags = []TagStats{}

	// Get total bookmarks
	h.db.Model(&models.Bookmark{}).Count(&stats.TotalBookmarks)

	// Get active bookmarks
	h.db.Model(&models.Bookmark{}).Where("archived = ?", false).Count(&stats.ActiveBookmarks)

	// Get archived bookmarks
	h.db.Model(&models.Bookmark{}).Where("archived = ?", true).Count(&stats.ArchivedBookmarks)

	// Get total tags
	h.db.Model(&models.Tag{}).Count(&stats.TotalTags)

	// First get special tags (To do and To read)
	specialRows, err := h.db.Raw(`
		SELECT 
			t.name,
			COUNT(DISTINCT bt.bookmark_id) as count,
			COUNT(DISTINCT CASE WHEN bt.completed THEN bt.bookmark_id END) as completed_count
		FROM tags t
		LEFT JOIN bookmark_tags bt ON bt.tag_id = t.id
		WHERE t.name IN ('To do', 'To read')
		GROUP BY t.id, t.name
	`).Rows()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer specialRows.Close()

	// Create a map to store special tags
	specialTags := make(map[string]TagStats)
	for specialRows.Next() {
		var tag TagStats
		if err := specialRows.Scan(&tag.Name, &tag.Count, &tag.CompletedCount); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		specialTags[tag.Name] = tag
	}

	// Add special tags first
	if todoTag, exists := specialTags["To do"]; exists {
		stats.TopTags = append(stats.TopTags, todoTag)
	}
	if toReadTag, exists := specialTags["To read"]; exists {
		stats.TopTags = append(stats.TopTags, toReadTag)
	}

	// Then get top tags (excluding To do and To read)
	rows, err := h.db.Raw(`
		SELECT 
			t.name,
			COUNT(DISTINCT bt.bookmark_id) as count,
			COUNT(DISTINCT CASE WHEN bt.completed THEN bt.bookmark_id END) as completed_count
		FROM tags t
		LEFT JOIN bookmark_tags bt ON bt.tag_id = t.id
		WHERE t.name NOT IN ('To do', 'To read')
		GROUP BY t.id, t.name
		HAVING COUNT(DISTINCT bt.bookmark_id) > 0
		ORDER BY count DESC
		LIMIT 5
	`).Rows()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	// Add regular top tags
	for rows.Next() {
		var tag TagStats
		if err := rows.Scan(&tag.Name, &tag.Count, &tag.CompletedCount); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		stats.TopTags = append(stats.TopTags, tag)
	}

	c.JSON(http.StatusOK, stats)
}

func (h *BookmarkHandler) ToggleArchive(c *gin.Context) {
	if err := h.service.ToggleArchive(c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Bookmark archive status toggled successfully"})
}

// Add this new handler method
func (h *BookmarkHandler) ToggleTagCompletion(c *gin.Context) {
	bookmarkID := c.Param("id")
	tagName := c.Param("tagName")

	log.Printf("ToggleTagCompletion: Starting for bookmarkID=%s, tagName=%s", bookmarkID, tagName)

	// Find the tag
	var tag models.Tag
	if err := h.db.Where("name = ?", tagName).First(&tag).Error; err != nil {
		log.Printf("ToggleTagCompletion: Tag not found: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Tag not found"})
		return
	}
	log.Printf("ToggleTagCompletion: Found tag with ID=%d", tag.ID)

	// Find the bookmark_tag entry
	var bookmarkTag models.BookmarkTag
	result := h.db.Where("bookmark_id = ? AND tag_id = ?", bookmarkID, tag.ID).
		First(&bookmarkTag)

	if result.Error != nil {
		log.Printf("ToggleTagCompletion: Creating new bookmark_tag entry")
		// If the bookmark_tag doesn't exist, create it
		bookmarkTag = models.BookmarkTag{
			BookmarkID: bookmarkID,
			TagID:      tag.ID,
			Completed:  true, // Start as completed since we're toggling
		}
		if err := h.db.Create(&bookmarkTag).Error; err != nil {
			log.Printf("ToggleTagCompletion: Error creating bookmark_tag: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		log.Printf("ToggleTagCompletion: Created new bookmark_tag with completed=true")
	} else {
		// Toggle the existing bookmark_tag's completed status
		newStatus := !bookmarkTag.Completed
		log.Printf("ToggleTagCompletion: Toggling existing bookmark_tag from %v to %v",
			bookmarkTag.Completed, newStatus)

		if err := h.db.Model(&bookmarkTag).Update("completed", newStatus).Error; err != nil {
			log.Printf("ToggleTagCompletion: Error updating bookmark_tag: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		bookmarkTag.Completed = newStatus
		log.Printf("ToggleTagCompletion: Updated bookmark_tag completed status")
	}

	log.Printf("ToggleTagCompletion: Returning completed=%v", bookmarkTag.Completed)
	c.JSON(http.StatusOK, gin.H{"completed": bookmarkTag.Completed})
}
