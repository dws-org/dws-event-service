package events

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

// Hilfsfunktion: Router nur für CreateEvent bauen
func setupRouterForCreate(controller *Controller) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/events", controller.CreateEvent)
	return r
}

func TestCreateEvent_InvalidJSON_Returns400(t *testing.T) {
	// Controller kann “leer” sein, weil wir bei invalid JSON nie zur DB kommen
	ec := &Controller{} // dbService wird hier nicht genutzt

	r := setupRouterForCreate(ec)

	req := httptest.NewRequest("POST", "/events", strings.NewReader(`{invalid json`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d. body=%s", w.Code, w.Body.String())
	}
}

func TestCreateEvent_MissingRequiredField_Returns400(t *testing.T) {
	ec := &Controller{}
	r := setupRouterForCreate(ec)

	// name fehlt absichtlich (binding:"required")
	body := `{
		"description":"desc",
		"startDate":"2026-01-01T00:00:00Z",
		"startTime":"2026-01-01T10:00:00Z",
		"price":"12.34",
		"endDate":"2026-01-02T00:00:00Z",
		"location":"Berlin",
		"capacity":10,
		"imageUrl":"https://example.com/a.jpg",
		"category":"workshop",
		"organizerId":"org-1"
	}`

	req := httptest.NewRequest("POST", "/events", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d. body=%s", w.Code, w.Body.String())
	}
}
