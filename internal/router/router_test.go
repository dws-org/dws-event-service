package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRouter_CORSMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	router := gin.New()
	
	// Add same CORS middleware as in production
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, x-api-key")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})
	
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	tests := []struct {
		name   string
		method string
		code   int
	}{
		{"options_request", http.MethodOptions, http.StatusNoContent},
		{"get_request", http.MethodGet, http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/test", nil)
			req.Header.Set("Origin", "http://example.com")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.code, w.Code)
			assert.Contains(t, w.Header().Get("Access-Control-Allow-Origin"), "*")
		})
	}
}

func TestRouter_HTTPMethods(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	router := gin.New()
	
	router.GET("/resource", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{})
	})
	router.POST("/resource", func(c *gin.Context) {
		c.JSON(http.StatusCreated, gin.H{})
	})
	router.PUT("/resource/:id", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{})
	})
	router.PATCH("/resource/:id", func(c *gin.Context) {
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
		{http.MethodPatch, "/resource/456", http.StatusOK},
		{http.MethodDelete, "/resource/789", http.StatusNoContent},
	}

	for _, tt := range tests {
		t.Run(tt.method+"_"+tt.path, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}

func TestRouter_PathParameters(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	router := gin.New()
	
	router.GET("/events/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.JSON(http.StatusOK, gin.H{"id": id})
	})
	
	router.GET("/events/:id/tickets/:ticketId", func(c *gin.Context) {
		eventId := c.Param("id")
		ticketId := c.Param("ticketId")
		c.JSON(http.StatusOK, gin.H{
			"eventId":  eventId,
			"ticketId": ticketId,
		})
	})

	tests := []struct {
		name       string
		path       string
		wantStatus int
	}{
		{"single_param", "/events/123", http.StatusOK},
		{"multiple_params", "/events/456/tickets/789", http.StatusOK},
		{"uuid_param", "/events/550e8400-e29b-41d4-a716-446655440000", http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestRouter_QueryParameters(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	router := gin.New()
	
	router.GET("/search", func(c *gin.Context) {
		query := c.Query("q")
		limit := c.DefaultQuery("limit", "10")
		c.JSON(http.StatusOK, gin.H{
			"query": query,
			"limit": limit,
		})
	})

	tests := []struct {
		name string
		path string
		want int
	}{
		{"with_query", "/search?q=test", http.StatusOK},
		{"with_multiple_params", "/search?q=golang&limit=20", http.StatusOK},
		{"no_query", "/search", http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.want, w.Code)
		})
	}
}

func TestRouter_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	router := gin.New()
	router.GET("/exists", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{})
	})

	req := httptest.NewRequest(http.MethodGet, "/does-not-exist", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouter_MethodNotAllowed(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	router := gin.New()
	router.GET("/resource", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{})
	})

	// Try POST on GET-only endpoint
	req := httptest.NewRequest(http.MethodPost, "/resource", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code) // Gin returns 404 for method not allowed
}
