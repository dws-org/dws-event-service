package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestErrorHandle_Middleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(ErrorHandle())
	
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestErrorHandle_WithError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(ErrorHandle())
	
	router.GET("/error", func(c *gin.Context) {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "something went wrong",
		})
	})

	req := httptest.NewRequest(http.MethodGet, "/error", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}

func TestErrorHandle_MultipleErrors(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(ErrorHandle())

	tests := []struct {
		name         string
		statusCode   int
		errorMessage string
	}{
		{"bad_request", http.StatusBadRequest, "Invalid input"},
		{"unauthorized", http.StatusUnauthorized, "Not authorized"},
		{"forbidden", http.StatusForbidden, "Access denied"},
		{"not_found", http.StatusNotFound, "Resource not found"},
		{"internal_error", http.StatusInternalServerError, "Server error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router.GET("/"+tt.name, func(c *gin.Context) {
				c.AbortWithStatusJSON(tt.statusCode, gin.H{"error": tt.errorMessage})
			})

			req := httptest.NewRequest(http.MethodGet, "/"+tt.name, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)
			assert.Contains(t, w.Body.String(), tt.errorMessage)
		})
	}
}

func TestErrorHandle_WithPanic(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(gin.Recovery()) // Add recovery middleware
	router.Use(ErrorHandle())
	
	router.GET("/panic", func(c *gin.Context) {
		panic("test panic")
	})

	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Recovery middleware should catch panic and return 500
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestMiddleware_ChainExecution(t *testing.T) {
	gin.SetMode(gin.TestMode)

	executionOrder := []string{}

	middleware1 := func(c *gin.Context) {
		executionOrder = append(executionOrder, "middleware1")
		c.Next()
	}

	middleware2 := func(c *gin.Context) {
		executionOrder = append(executionOrder, "middleware2")
		c.Next()
	}

	router := gin.New()
	router.Use(middleware1, middleware2, ErrorHandle())
	
	router.GET("/test", func(c *gin.Context) {
		executionOrder = append(executionOrder, "handler")
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, []string{"middleware1", "middleware2", "handler"}, executionOrder)
}

func TestMiddleware_AbortChain(t *testing.T) {
	gin.SetMode(gin.TestMode)

	middlewareAbort := func(c *gin.Context) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	}

	router := gin.New()
	router.Use(ErrorHandle(), middlewareAbort)
	
	router.GET("/test", func(c *gin.Context) {
		t.Error("Handler should not be called")
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestMiddleware_ContextValues(t *testing.T) {
	gin.SetMode(gin.TestMode)

	setValueMiddleware := func(c *gin.Context) {
		c.Set("user_id", "123")
		c.Set("request_id", "req-456")
		c.Next()
	}

	router := gin.New()
	router.Use(setValueMiddleware)
	
	router.GET("/test", func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		assert.True(t, exists)
		assert.Equal(t, "123", userID)

		requestID, exists := c.Get("request_id")
		assert.True(t, exists)
		assert.Equal(t, "req-456", requestID)

		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestMiddleware_HeaderManipulation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	addHeadersMiddleware := func(c *gin.Context) {
		c.Header("X-Custom-Header", "custom-value")
		c.Header("X-Request-ID", "12345")
		c.Next()
	}

	router := gin.New()
	router.Use(addHeadersMiddleware)
	
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, "custom-value", w.Header().Get("X-Custom-Header"))
	assert.Equal(t, "12345", w.Header().Get("X-Request-ID"))
}
