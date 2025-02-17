package jobs

import (
	"log"
	"time"

	"github.com/helioLJ/tweetvault/internal/services"
)

func StartViewRefreshJob(bookmarkService *services.BookmarkService) {
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := bookmarkService.RefreshView(); err != nil {
					log.Printf("Error refreshing materialized view: %v", err)
				}
			}
		}
	}()
}
