package health

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

// TestLive tests the liveness probe endpoint
func TestLive(t *testing.T) {
	router := gin.New()
	controller := NewController()
	router.GET("/livez", controller.Live)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/livez", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "status")
}

// TestReady tests the readiness probe endpoint
func TestReady(t *testing.T) {
	router := gin.New()
	controller := NewController()
	router.GET("/readyz", controller.Ready)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/readyz", nil)
	router.ServeHTTP(w, req)

	// Should return 200 or 503 depending on database connection
	assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusServiceUnavailable)
}

// TestInfo tests the info endpoint
func TestInfo(t *testing.T) {
router := gin.New()
controller := NewController()
router.GET("/_meta", controller.Info)

w := httptest.NewRecorder()
req, _ := http.NewRequest("GET", "/_meta", nil)
router.ServeHTTP(w, req)

assert.Equal(t, http.StatusOK, w.Code)
assert.Contains(t, w.Body.String(), "service")
}

// TestNewController tests controller initialization
func TestNewController(t *testing.T) {
	controller := NewController()
	assert.NotNil(t, controller)
	assert.NotNil(t, controller.service)
}
