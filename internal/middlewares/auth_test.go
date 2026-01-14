package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Test-Modus für Gin (keine Debug-Ausgaben)
func init() {
	gin.SetMode(gin.TestMode)
}

// =============================================================================
// TEST: AuthMiddleware
// =============================================================================

func TestAuthMiddleware_OhneToken(t *testing.T) {
	// Setup: Router mit AuthMiddleware
	router := gin.New()
	router.Use(AuthMiddleware())
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Request OHNE Authorization Header
	req, _ := http.NewRequest("GET", "/protected", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	// Erwartung: 401 Unauthorized
	assert.Equal(t, http.StatusUnauthorized, resp.Code)
}

func TestAuthMiddleware_MitGueltigemToken(t *testing.T) {
	// Setup: Router mit AuthMiddleware
	router := gin.New()
	router.Use(AuthMiddleware())
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Request MIT gültigem Bearer Token
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer mein-test-token")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	// Erwartung: 200 OK
	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestAuthMiddleware_FalschesFormat(t *testing.T) {
	router := gin.New()
	router.Use(AuthMiddleware())
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Request mit falschem Format (ohne "Bearer")
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "mein-token-ohne-bearer")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	// Erwartung: 401 Unauthorized
	assert.Equal(t, http.StatusUnauthorized, resp.Code)
}

func TestAuthMiddleware_LeeresToken(t *testing.T) {
	router := gin.New()
	router.Use(AuthMiddleware())
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Request mit leerem Token
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer ")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	// Erwartung: 401 Unauthorized
	assert.Equal(t, http.StatusUnauthorized, resp.Code)
}

// =============================================================================
// TEST: GetUserIDFromContext
// =============================================================================

func TestGetUserIDFromContext(t *testing.T) {
	router := gin.New()
	router.Use(AuthMiddleware())
	router.GET("/user", func(c *gin.Context) {
		userID, exists := GetUserIDFromContext(c)
		c.JSON(http.StatusOK, gin.H{
			"userID": userID,
			"exists": exists,
		})
	})

	req, _ := http.NewRequest("GET", "/user", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "userID")
}