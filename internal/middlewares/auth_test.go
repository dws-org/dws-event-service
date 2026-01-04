package middlewares

import (
"net/http"
"net/http/httptest"
"testing"

"github.com/gin-gonic/gin"
"github.com/stretchr/testify/assert"
)

func init() {
gin.SetMode(gin.TestMode)
}

// TestKeycloakAuthMiddleware_MissingToken tests 401 error when no token provided
// REQ7: Failure test case - 401 Unauthorized (missing token)
func TestKeycloakAuthMiddleware_MissingToken(t *testing.T) {
router := gin.New()
router.Use(KeycloakAuthMiddleware())
router.GET("/protected", func(c *gin.Context) {
c.JSON(http.StatusOK, gin.H{"message": "success"})
})

w := httptest.NewRecorder()
req, _ := http.NewRequest("GET", "/protected", nil)
// No Authorization header
router.ServeHTTP(w, req)

assert.Equal(t, http.StatusUnauthorized, w.Code, "Expected 401 Unauthorized")
assert.Contains(t, w.Body.String(), "Authorization header missing")
}

// TestKeycloakAuthMiddleware_InvalidTokenFormat tests 401 error with malformed token
// REQ7: Failure test case - 401 Unauthorized (invalid token format)
func TestKeycloakAuthMiddleware_InvalidTokenFormat(t *testing.T) {
router := gin.New()
router.Use(KeycloakAuthMiddleware())
router.GET("/protected", func(c *gin.Context) {
c.JSON(http.StatusOK, gin.H{"message": "success"})
})

w := httptest.NewRecorder()
req, _ := http.NewRequest("GET", "/protected", nil)
req.Header.Set("Authorization", "InvalidFormat token123")
router.ServeHTTP(w, req)

assert.Equal(t, http.StatusUnauthorized, w.Code, "Expected 401 Unauthorized")
}

// TestKeycloakAuthMiddleware_EmptyToken tests 401 error with empty Bearer token
// REQ7: Failure test case - 401 Unauthorized (empty token)
func TestKeycloakAuthMiddleware_EmptyToken(t *testing.T) {
router := gin.New()
router.Use(KeycloakAuthMiddleware())
router.GET("/protected", func(c *gin.Context) {
c.JSON(http.StatusOK, gin.H{"message": "success"})
})

w := httptest.NewRecorder()
req, _ := http.NewRequest("GET", "/protected", nil)
req.Header.Set("Authorization", "Bearer ")
router.ServeHTTP(w, req)

assert.Equal(t, http.StatusUnauthorized, w.Code, "Expected 401 Unauthorized")
}

// TestErrorHandler tests error handling middleware
func TestErrorHandler(t *testing.T) {
router := gin.New()
router.Use(ErrorHandler())
router.GET("/error", func(c *gin.Context) {
c.Error(assert.AnError)
c.JSON(http.StatusInternalServerError, gin.H{"error": "test error"})
})

w := httptest.NewRecorder()
req, _ := http.NewRequest("GET", "/error", nil)
router.ServeHTTP(w, req)

assert.Equal(t, http.StatusInternalServerError, w.Code)
}
