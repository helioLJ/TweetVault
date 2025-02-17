package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/helioLJ/tweetvault/config"
	"github.com/helioLJ/tweetvault/internal/api/routes"
	"github.com/helioLJ/tweetvault/internal/database"
	"github.com/helioLJ/tweetvault/internal/jobs"
	"github.com/helioLJ/tweetvault/internal/services"
)

func main() {
	// Enable Gin debug mode
	gin.SetMode(gin.DebugMode)

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database connection
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Create bookmark service
	bookmarkService := services.NewBookmarkService(db)

	// Initialize router with custom logging
	r := routes.SetupRouter(db)

	// Add custom logging middleware
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// Return detailed error messages if present
		if param.ErrorMessage != "" {
			return "ERROR: " + param.ErrorMessage + "\n"
		}
		return ""
	}))

	jobs.StartViewRefreshJob(bookmarkService)

	// Start server
	log.Printf("Server starting on port %s", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
