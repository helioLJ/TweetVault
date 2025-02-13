package routes

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/helioLJ/tweetvault/internal/api/handlers"
	"github.com/helioLJ/tweetvault/internal/api/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.New()

	// Add recovery middleware
	r.Use(gin.Recovery())

	// Add custom logger middleware
	r.Use(gin.LoggerWithWriter(gin.DefaultWriter, "/api/health"))

	// Add error logging middleware
	r.Use(func(c *gin.Context) {
		c.Next()

		// Print any errors that occurred
		if len(c.Errors) > 0 {
			log.Printf("Errors: %v", c.Errors)
		}
	})

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

	// Add metrics middleware to all routes
	r.Use(middleware.MetricsMiddleware)

	// Add Prometheus metrics endpoint
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

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
		api.POST("/bookmarks/:id/toggle-archive", bookmarkHandler.ToggleArchive)
		api.POST("/bookmarks/:id/tags/:tagName/toggle-completion", bookmarkHandler.ToggleTagCompletion)

		// Tag endpoints
		api.GET("/tags", tagHandler.List)
		api.POST("/tags", tagHandler.Create)
		api.PUT("/tags/:id", tagHandler.Update)
		api.DELETE("/tags/:id", tagHandler.Delete)
		api.GET("/tags/:id/count", tagHandler.GetBookmarkCount)

		// Statistics endpoint
		api.GET("/statistics", bookmarkHandler.GetStatistics)

		// Metrics endpoints
		api.POST("/metrics/error", handlers.HandleErrorMetric)
		api.POST("/metrics/timing", handlers.HandleTimingMetric)
		api.POST("/metrics/pageview", handlers.HandlePageViewMetric)
		api.POST("/metrics/performance", handlers.HandlePerformanceMetric)
		api.POST("/metrics/vitals", handlers.HandleWebVitalsMetric)
		api.POST("/metrics/bookmark-load", handlers.HandleBookmarkLoadMetric)
	}

	return r
}
