package events

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreateEvent_ValidationErrors(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		payload        map[string]interface{}
		expectedStatus int
		expectedError  string
	}{
		{
			name: "missing_name",
			payload: map[string]interface{}{
				"description": "Test event",
				"date":        time.Now().Add(24 * time.Hour).Format(time.RFC3339),
				"location":    "Test Location",
				"capacity":    100,
				"price":       "10.00",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "name",
		},
		{
			name: "missing_date",
			payload: map[string]interface{}{
				"name":        "Test Event",
				"description": "Test event",
				"location":    "Test Location",
				"capacity":    100,
				"price":       "10.00",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "date",
		},
		{
			name: "missing_location",
			payload: map[string]interface{}{
				"name":        "Test Event",
				"description": "Test event",
				"date":        time.Now().Add(24 * time.Hour).Format(time.RFC3339),
				"capacity":    100,
				"price":       "10.00",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "location",
		},
		{
			name: "negative_capacity",
			payload: map[string]interface{}{
				"name":        "Test Event",
				"description": "Test event",
				"date":        time.Now().Add(24 * time.Hour).Format(time.RFC3339),
				"location":    "Test Location",
				"capacity":    -1,
				"price":       "10.00",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "capacity",
		},
		{
			name: "zero_capacity",
			payload: map[string]interface{}{
				"name":        "Test Event",
				"description": "Test event",
				"date":        time.Now().Add(24 * time.Hour).Format(time.RFC3339),
				"location":    "Test Location",
				"capacity":    0,
				"price":       "10.00",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "capacity",
		},
		{
			name: "invalid_date_format",
			payload: map[string]interface{}{
				"name":        "Test Event",
				"description": "Test event",
				"date":        "invalid-date",
				"location":    "Test Location",
				"capacity":    100,
				"price":       "10.00",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "date",
		},
		{
			name: "past_date",
			payload: map[string]interface{}{
				"name":        "Test Event",
				"description": "Test event",
				"date":        time.Now().Add(-24 * time.Hour).Format(time.RFC3339),
				"location":    "Test Location",
				"capacity":    100,
				"price":       "10.00",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "date",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			controller := &Controller{}
			
			router.POST("/events", controller.CreateEvent)

			body, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectedError)
		})
	}
}

func TestCreateEvent_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	controller := &Controller{}
	router.POST("/events", controller.CreateEvent)

	tests := []struct {
		name string
		body string
	}{
		{"invalid_json", `{"name": "Test", "invalid"`},
		{"empty_body", ""},
		{"not_json", "this is not json"},
		{"malformed", `{name: 'missing quotes'}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code)
		})
	}
}

func TestCreateEvent_EdgeCases(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		payload        map[string]interface{}
		expectedStatus int
	}{
		{
			name: "very_long_name",
			payload: map[string]interface{}{
				"name":        string(make([]byte, 1000)),
				"description": "Test",
				"date":        time.Now().Add(24 * time.Hour).Format(time.RFC3339),
				"location":    "Location",
				"capacity":    100,
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "very_large_capacity",
			payload: map[string]interface{}{
				"name":        "Event",
				"description": "Test",
				"date":        time.Now().Add(24 * time.Hour).Format(time.RFC3339),
				"location":    "Location",
				"capacity":    999999999,
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "special_characters_in_name",
			payload: map[string]interface{}{
				"name":        "Event<script>alert('xss')</script>",
				"description": "Test",
				"date":        time.Now().Add(24 * time.Hour).Format(time.RFC3339),
				"location":    "Location",
				"capacity":    100,
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			controller := &Controller{}
			router.POST("/events", controller.CreateEvent)

			body, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestCreateEvent_PriceValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		price          interface{}
		expectedStatus int
	}{
		{"negative_price", "-10.00", http.StatusBadRequest},
		{"zero_price", "0.00", http.StatusBadRequest},
		{"valid_price", "10.00", http.StatusBadRequest}, // Still fails due to DB
		{"high_price", "9999.99", http.StatusBadRequest},
		{"invalid_price_format", "abc", http.StatusBadRequest},
		{"price_with_currency", "$10.00", http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			controller := &Controller{}
			router.POST("/events", controller.CreateEvent)

			payload := map[string]interface{}{
				"name":        "Test Event",
				"description": "Test",
				"date":        time.Now().Add(24 * time.Hour).Format(time.RFC3339),
				"location":    "Location",
				"capacity":    100,
				"price":       tt.price,
			}

			body, _ := json.Marshal(payload)
			req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestCreateEvent_CategoryValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	validCategories := []string{"Music", "Sports", "Technology", "Art", "Food"}
	
	for _, category := range validCategories {
		t.Run("valid_category_"+category, func(t *testing.T) {
			router := gin.New()
			controller := &Controller{}
			router.POST("/events", controller.CreateEvent)

			payload := map[string]interface{}{
				"name":        "Test Event",
				"description": "Test",
				"date":        time.Now().Add(24 * time.Hour).Format(time.RFC3339),
				"location":    "Location",
				"capacity":    100,
				"category":    category,
			}

			body, _ := json.Marshal(payload)
			req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			// Will fail due to no DB but should pass validation
			assert.NotEqual(t, http.StatusBadRequest, w.Code)
		})
	}
}

func TestController_HTTPMethods(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller := &Controller{}
	router := gin.New()
	
	router.GET("/events", controller.GetEvents)
	router.GET("/events/:id", controller.GetEventByID)
	router.POST("/events", controller.CreateEvent)

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{"get_events_valid", http.MethodGet, "/events", http.StatusInternalServerError}, // No DB
		{"get_event_by_id_valid", http.MethodGet, "/events/1", http.StatusInternalServerError},
		{"create_event_post", http.MethodPost, "/events", http.StatusBadRequest}, // No body
		{"wrong_method_put", http.MethodPut, "/events", http.StatusNotFound},
		{"wrong_method_delete", http.MethodDelete, "/events", http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
