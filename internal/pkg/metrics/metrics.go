package metrics

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTP request counter
	RequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status", "service"},
	)

	// HTTP request duration histogram
	RequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request latency in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint", "service"},
	)

	// Event operations counter
	EventOperations = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "event_operations_total",
			Help: "Total number of event operations",
		},
		[]string{"operation", "status", "service"},
	)

	// RabbitMQ messages counter
	RabbitMQMessages = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "rabbitmq_messages_total",
			Help: "Total number of RabbitMQ messages",
		},
		[]string{"action", "queue", "status", "service"},
	)

	// Database operations counter
	DatabaseOperations = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "database_operations_total",
			Help: "Total number of database operations",
		},
		[]string{"operation", "table", "status", "service"},
	)
)

// PrometheusMiddleware records HTTP metrics
func PrometheusMiddleware(serviceName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		duration := time.Since(start).Seconds()
		status := c.Writer.Status()

		// Record metrics
		RequestsTotal.WithLabelValues(
			c.Request.Method,
			c.FullPath(),
			fmt.Sprintf("%d", status),
			serviceName,
		).Inc()

		RequestDuration.WithLabelValues(
			c.Request.Method,
			c.FullPath(),
			serviceName,
		).Observe(duration)
	}
}
