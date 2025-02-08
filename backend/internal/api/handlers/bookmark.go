package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"github.com/helioLJ/tweetvault/internal/models"
	"github.com/helioLJ/tweetvault/internal/services"
)

type BookmarkHandler struct {
	service *services.BookmarkService
}

func NewBookmarkHandler(db *gorm.DB) *BookmarkHandler {
	return &BookmarkHandler{
		service: services.NewBookmarkService(db),
	}
}

// List returns all bookmarks with optional filtering
func (h *BookmarkHandler) List(c *gin.Context) {
	bookmarks, total, err := h.service.List(c.Query("tag"))
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
