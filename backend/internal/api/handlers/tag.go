package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/helioLJ/tweetvault/internal/models"
	"gorm.io/gorm"
)

type TagHandler struct {
	db *gorm.DB
}

var standardTags = []string{"To do", "To read"}

func NewTagHandler(db *gorm.DB) *TagHandler {
	return &TagHandler{db: db}
}

func (h *TagHandler) List(c *gin.Context) {
	var tags []models.Tag
	if err := h.db.Find(&tags).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tags)
}

func (h *TagHandler) Update(c *gin.Context) {
	var input struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var tag models.Tag
	if err := h.db.First(&tag, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tag not found"})
		return
	}

	if isStandardTag(tag.Name) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot rename standard tags"})
		return
	}

	tag.Name = input.Name
	if err := h.db.Save(&tag).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tag)
}

func (h *TagHandler) Delete(c *gin.Context) {
	var tag models.Tag
	if err := h.db.First(&tag, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tag not found"})
		return
	}

	if isStandardTag(tag.Name) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot delete standard tags"})
		return
	}

	// Start a transaction
	tx := h.db.Begin()

	// Delete the tag associations from bookmark_tags
	if err := tx.Exec("DELETE FROM bookmark_tags WHERE tag_id = ?", tag.ID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Delete the tag itself
	if err := tx.Delete(&tag).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update tag count after successful deletion
	var totalTags int64
	h.db.Model(&models.Tag{}).Count(&totalTags)

	c.JSON(http.StatusOK, gin.H{"message": "Tag deleted successfully"})
}

func (h *TagHandler) GetBookmarkCount(c *gin.Context) {
	var count int64
	if err := h.db.Model(&models.Bookmark{}).
		Joins("JOIN bookmark_tags ON bookmarks.id = bookmark_tags.bookmark_id").
		Where("bookmark_tags.tag_id = ?", c.Param("id")).
		Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": count})
}

func (h *TagHandler) Create(c *gin.Context) {
	var input struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tag := models.Tag{Name: input.Name}
	if err := h.db.Create(&tag).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update tag count after successful creation
	var totalTags int64
	h.db.Model(&models.Tag{}).Count(&totalTags)

	c.JSON(http.StatusCreated, tag)
}

func isStandardTag(tagName string) bool {
	for _, st := range standardTags {
		if st == tagName {
			return true
		}
	}
	return false
}
