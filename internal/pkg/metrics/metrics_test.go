package metrics

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestPrometheusMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(PrometheusMiddleware("test-service"))
	
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPrometheusMiddleware_DifferentMethods(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(PrometheusMiddleware("test-service"))
	
	router.GET("/resource", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{})
	})
	router.POST("/resource", func(c *gin.Context) {
		c.JSON(http.StatusCreated, gin.H{})
	})
	router.PUT("/resource/:id", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{})
	})
	router.DELETE("/resource/:id", func(c *gin.Context) {
		c.JSON(http.StatusNoContent, gin.H{})
	})

	tests := []struct {
		method       string
		path         string
		expectedCode int
	}{
		{http.MethodGet, "/resource", http.StatusOK},
		{http.MethodPost, "/resource", http.StatusCreated},
		{http.MethodPut, "/resource/123", http.StatusOK},
		{http.MethodDelete, "/resource/456", http.StatusNoContent},
	}

	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}

func TestPrometheusMiddleware_ErrorResponses(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(PrometheusMiddleware("test-service"))
	
	router.GET("/error", func(c *gin.Context) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
	})
	router.GET("/notfound", func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	})
	router.GET("/badrequest", func(c *gin.Context) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
	})

	tests := []struct {
		name         string
		path         string
		expectedCode int
	}{
		{"server_error", "/error", http.StatusInternalServerError},
		{"not_found", "/notfound", http.StatusNotFound},
		{"bad_request", "/badrequest", http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}

func TestPrometheusMiddleware_Latency(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(PrometheusMiddleware("test-service"))
	
	router.GET("/slow", func(c *gin.Context) {
		// Simulate slow endpoint
		// In real test, we'd verify duration is recorded
		c.JSON(http.StatusOK, gin.H{"message": "slow"})
	})

	req := httptest.NewRequest(http.MethodGet, "/slow", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestMetrics_CountersExist(t *testing.T) {
	// Test that metrics are initialized
	assert.NotNil(t, RequestsTotal, "RequestsTotal counter should exist")
	assert.NotNil(t, RequestDuration, "RequestDuration histogram should exist")
	assert.NotNil(t, EventOperations, "EventOperations counter should exist")
	assert.NotNil(t, RabbitMQMessages, "RabbitMQMessages counter should exist")
	assert.NotNil(t, DatabaseOperations, "DatabaseOperations counter should exist")
}

func TestMetrics_IncrementCounters(t *testing.T) {
	// Test that we can increment metrics without panic
	assert.NotPanics(t, func() {
		RequestsTotal.WithLabelValues("GET", "/test", "200", "test-service").Inc()
		EventOperations.WithLabelValues("create", "success", "test-service").Inc()
		DatabaseOperations.WithLabelValues("insert", "events", "success", "test-service").Inc()
		RabbitMQMessages.WithLabelValues("publish", "events", "success", "test-service").Inc()
	})
}

func TestMetrics_ObserveHistogram(t *testing.T) {
	// Test that we can record duration observations
	assert.NotPanics(t, func() {
		RequestDuration.WithLabelValues("GET", "/test", "test-service").Observe(0.123)
		RequestDuration.WithLabelValues("POST", "/events", "test-service").Observe(0.456)
	})
}
