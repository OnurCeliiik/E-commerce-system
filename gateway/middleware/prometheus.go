package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// Total number of HTTP requests
	httpRequestTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	// HTTP request duration in seconds
	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "status"},
	)
)

func init() {
	prometheus.MustRegister(httpRequestTotal)
	prometheus.MustRegister(httpRequestDuration)
}

func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Request start time
		start := time.Now()
		// Send it to the handler
		c.Next()

		// Metric collection after the request is handled
		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())
		method := c.Request.Method
		path := c.FullPath()

		// If path is empty, set it to "unknown"
		if path == "" {
			path = "unknown"
		}

		// Update metrics
		httpRequestTotal.WithLabelValues(method, path, status).Inc()
		httpRequestDuration.WithLabelValues(method, path, status).Observe(duration)
	}
}
