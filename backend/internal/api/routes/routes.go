package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/helioLJ/tweetvault/internal/api/handlers"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// Enable CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})

	// Create handler instances
	uploadHandler := handlers.NewUploadHandler(db)
	bookmarkHandler := handlers.NewBookmarkHandler(db)
	tagHandler := handlers.NewTagHandler(db)

	// API routes
	api := r.Group("/api")
	{
		// Upload endpoints
		api.POST("/upload", uploadHandler.HandleUpload)

		// Bookmark endpoints
		api.GET("/bookmarks", bookmarkHandler.List)
		api.GET("/bookmarks/:id", bookmarkHandler.Get)
		api.PUT("/bookmarks/:id", bookmarkHandler.Update)
		api.DELETE("/bookmarks/:id", bookmarkHandler.Delete)

		// Tag endpoints
		api.GET("/tags", tagHandler.List)
	}

	return r
}