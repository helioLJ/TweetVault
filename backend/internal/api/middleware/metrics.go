package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests",
			Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
		},
		[]string{"method", "endpoint"},
	)
)

func MetricsMiddleware(c *gin.Context) {
	start := time.Now()

	c.Next()

	duration := time.Since(start).Seconds()
	status := c.Writer.Status()

	httpRequestsTotal.WithLabelValues(c.Request.Method, c.Request.URL.Path, string(rune(status))).Inc()
	httpRequestDuration.WithLabelValues(c.Request.Method, c.Request.URL.Path).Observe(duration)
}

// ResponseWriterWrapper to capture status code
type ResponseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

func NewResponseWriterWrapper(w http.ResponseWriter) *ResponseWriterWrapper {
	return &ResponseWriterWrapper{w, http.StatusOK}
}

func (wrw *ResponseWriterWrapper) WriteHeader(code int) {
	wrw.statusCode = code
	wrw.ResponseWriter.WriteHeader(code)
}
