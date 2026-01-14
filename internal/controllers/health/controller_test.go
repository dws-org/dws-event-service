package health

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/oskargbc/dws-event-service.git/configs"
)

// newTestController erstellt einen Health-Controller ohne configs.GetEnvConfig() aufzurufen.
// Dadurch vermeiden wir, dass beim Testlauf eine config.yaml vorhanden sein muss.
func newTestController() *Controller {
	return &Controller{
		service: configs.Service{
			Name:        "Test Service",
			Slug:        "test-service",
			Description: "Test description",
			Version:     "0.0.0-test",
			Tags:        []string{"test"},
		},
		// startedAt in der Vergangenheit, damit uptime nicht "0s" ist (optional).
		startedAt: time.Now().UTC().Add(-10 * time.Second),
	}
}

func setupHealthRouter(hc *Controller) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/health/live", hc.Live)
	r.GET("/health/info", hc.Info)
	return r
}

func TestLive_Returns200AndStatusOk(t *testing.T) {
	hc := newTestController()
	r := setupHealthRouter(hc)

	req := httptest.NewRequest(http.MethodGet, "/health/live", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d. body=%s", w.Code, w.Body.String())
	}

	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid json: %v, body=%s", err, w.Body.String())
	}

	// Pflichtfelder (Werte nicht exakt prüfen, weil timestamp/uptime dynamisch sind)
	if body["status"] != "ok" {
		t.Fatalf("expected status=ok, got %v", body["status"])
	}
	if _, ok := body["service"]; !ok {
		t.Fatalf("expected service field")
	}
	if _, ok := body["timestamp"]; !ok {
		t.Fatalf("expected timestamp field")
	}
	if _, ok := body["uptime"]; !ok {
		t.Fatalf("expected uptime field")
	}
}

func TestInfo_Returns200AndMetadataFields(t *testing.T) {
	hc := newTestController()
	r := setupHealthRouter(hc)

	req := httptest.NewRequest(http.MethodGet, "/health/info", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d. body=%s", w.Code, w.Body.String())
	}

	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid json: %v, body=%s", err, w.Body.String())
	}

	// Wir prüfen, dass die wichtigsten Metadaten-Felder vorhanden sind
	for _, k := range []string{"name", "slug", "description", "version", "tags", "startedAt", "uptime"} {
		if _, ok := body[k]; !ok {
			t.Fatalf("expected field %q", k)
		}
	}
}

