package handlers

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// Frontend error metrics
	frontendErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "frontend_errors_total",
			Help: "Total number of frontend errors",
		},
		[]string{"endpoint", "error_type"},
	)

	// API timing metrics
	apiTimings = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "api_request_duration_seconds",
			Help:    "API request duration in seconds",
			Buckets: []float64{0.1, 0.5, 1, 2, 5},
		},
		[]string{"endpoint"},
	)

	// Page view metrics
	pageViews = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "page_views_total",
			Help: "Total number of page views",
		},
		[]string{"page"},
	)

	// Page load timing metrics
	pageLoadTiming = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "page_load_duration_seconds",
			Help:    "Page load duration in seconds",
			Buckets: []float64{0.1, 0.5, 1, 2, 5},
		},
		[]string{"page"},
	)

	// Web vitals metrics
	webVitals = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "web_vitals",
			Help:    "Web Vitals metrics",
			Buckets: []float64{0.1, 0.5, 1, 2, 5},
		},
		[]string{"name"},
	)

	// Bookmark load metrics
	bookmarkLoadTiming = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "bookmark_load_duration_seconds",
			Help:    "Bookmark load duration in seconds",
			Buckets: []float64{0.1, 0.5, 1, 2, 5},
		},
		[]string{"status"},
	)

	bookmarkLoadCount = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "bookmark_load_count",
			Help:    "Number of bookmarks loaded per request",
			Buckets: []float64{10, 20, 50, 100, 200},
		},
		[]string{"status"},
	)
)

type ErrorMetric struct {
	Endpoint  string    `json:"endpoint"`
	Message   string    `json:"message"`
	Stack     string    `json:"stack"`
	Timestamp time.Time `json:"timestamp"`
}

type TimingMetric struct {
	Endpoint  string    `json:"endpoint"`
	Duration  float64   `json:"duration"`
	Timestamp time.Time `json:"timestamp"`
}

type PageViewMetric struct {
	Page      string    `json:"page"`
	LoadTime  float64   `json:"loadTime"`
	Timestamp time.Time `json:"timestamp"`
}

type PerformanceMetric struct {
	Page      string    `json:"page"`
	Metric    string    `json:"metric"`
	Value     float64   `json:"value"`
	Timestamp time.Time `json:"timestamp"`
}

type WebVitalsMetric struct {
	Name      string    `json:"name"`
	Value     float64   `json:"value"`
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
}

type BookmarkLoadMetric struct {
	Duration  float64   `json:"duration"`
	Count     int       `json:"count"`
	Timestamp time.Time `json:"timestamp"`
}

func HandleErrorMetric(c *gin.Context) {
	var metric ErrorMetric
	if err := c.BindJSON(&metric); err != nil {
		c.JSON(400, gin.H{"error": "Invalid metric data"})
		return
	}

	frontendErrors.WithLabelValues(metric.Endpoint, "error").Inc()
	c.Status(200)
}

func HandleTimingMetric(c *gin.Context) {
	var metric TimingMetric
	if err := c.BindJSON(&metric); err != nil {
		c.JSON(400, gin.H{"error": "Invalid metric data"})
		return
	}

	apiTimings.WithLabelValues(metric.Endpoint).Observe(metric.Duration / 1000.0) // Convert to seconds
	c.Status(200)
}

func HandlePageViewMetric(c *gin.Context) {
	var metric PageViewMetric
	if err := c.BindJSON(&metric); err != nil {
		c.JSON(400, gin.H{"error": "Invalid metric data"})
		return
	}

	pageViews.WithLabelValues(metric.Page).Inc()
	pageLoadTiming.WithLabelValues(metric.Page).Observe(metric.LoadTime / 1000.0)
	c.Status(200)
}

func HandlePerformanceMetric(c *gin.Context) {
	var metric PerformanceMetric
	if err := c.BindJSON(&metric); err != nil {
		c.JSON(400, gin.H{"error": "Invalid metric data"})
		return
	}

	// Record performance metrics based on the metric type
	webVitals.WithLabelValues(metric.Metric).Observe(metric.Value)
	c.Status(200)
}

func HandleWebVitalsMetric(c *gin.Context) {
	var metric WebVitalsMetric
	if err := c.BindJSON(&metric); err != nil {
		c.JSON(400, gin.H{"error": "Invalid metric data"})
		return
	}

	webVitals.WithLabelValues(metric.Name).Observe(metric.Value)
	c.Status(200)
}

func HandleBookmarkLoadMetric(c *gin.Context) {
	var metric BookmarkLoadMetric
	if err := c.BindJSON(&metric); err != nil {
		c.JSON(400, gin.H{"error": "Invalid metric data"})
		return
	}

	status := "success"
	if metric.Count == 0 {
		status = "empty"
	}

	bookmarkLoadTiming.WithLabelValues(status).Observe(metric.Duration / 1000.0)
	bookmarkLoadCount.WithLabelValues(status).Observe(float64(metric.Count))
	c.Status(200)
}
