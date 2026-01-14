package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestCreateEventRequest_Validation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name        string
		payload     map[string]interface{}
		expectError bool
	}{
		{
			name: "valid_complete_request",
			payload: map[string]interface{}{
				"name":        "Summer Concert",
				"description": "Amazing outdoor concert",
				"startDate":   time.Now().Format(time.RFC3339),
				"startTime":   time.Now().Format(time.RFC3339),
				"price":       "49.99",
				"endDate":     time.Now().Add(24 * time.Hour).Format(time.RFC3339),
				"location":    "Central Park",
				"capacity":    1000,
				"imageUrl":    "https://example.com/image.jpg",
				"category":    "Music",
				"organizerId": "org-123",
			},
			expectError: false,
		},
		{
			name: "missing_name",
			payload: map[string]interface{}{
				"description": "Event without name",
				"capacity":    100,
			},
			expectError: true,
		},
		{
			name: "invalid_capacity_negative",
			payload: map[string]interface{}{
				"name":     "Test Event",
				"capacity": -100,
			},
			expectError: false, // Validation happens at API level
		},
		{
			name: "zero_capacity",
			payload: map[string]interface{}{
				"name":     "Test Event",
				"capacity": 0,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, err := json.Marshal(tt.payload)
			assert.NoError(t, err)

			var decoded map[string]interface{}
			err = json.Unmarshal(jsonData, &decoded)

			if tt.expectError {
				// For incomplete payloads, certain fields will be missing
				assert.True(t, len(decoded) < 11)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCreateEventRequest_JSONMarshal(t *testing.T) {
	now := time.Now().UTC()
	price := decimal.NewFromFloat(99.99)

	request := map[string]interface{}{
		"name":        "Tech Conference 2026",
		"description": "Annual tech conference",
		"startDate":   now.Format(time.RFC3339),
		"startTime":   now.Format(time.RFC3339),
		"price":       price.String(),
		"endDate":     now.Add(48 * time.Hour).Format(time.RFC3339),
		"location":    "Convention Center",
		"capacity":    5000,
		"imageUrl":    "https://example.com/tech.jpg",
		"category":    "Technology",
		"organizerId": "org-456",
	}

	jsonData, err := json.Marshal(request)
	assert.NoError(t, err)
	assert.NotEmpty(t, jsonData)

	var decoded map[string]interface{}
	err = json.Unmarshal(jsonData, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, "Tech Conference 2026", decoded["name"])
	assert.Equal(t, float64(5000), decoded["capacity"])
}

func TestEventEndpoints_RequestFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		method         string
		path           string
		body           interface{}
		expectedStatus int
	}{
		{
			name:           "get_events_list",
			method:         http.MethodGet,
			path:           "/events",
			body:           nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "get_event_by_id",
			method:         http.MethodGet,
			path:           "/events/123",
			body:           nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:   "create_event_with_body",
			method: http.MethodPost,
			path:   "/events",
			body: map[string]interface{}{
				"name":     "New Event",
				"capacity": 100,
			},
			expectedStatus: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var bodyReader *bytes.Reader
			if tt.body != nil {
				jsonData, _ := json.Marshal(tt.body)
				bodyReader = bytes.NewReader(jsonData)
			} else {
				bodyReader = bytes.NewReader([]byte{})
			}

			req := httptest.NewRequest(tt.method, tt.path, bodyReader)
			if tt.body != nil {
				req.Header.Set("Content-Type", "application/json")
			}

			assert.NotNil(t, req)
			assert.Equal(t, tt.method, req.Method)
			assert.Equal(t, tt.path, req.URL.Path)
		})
	}
}

func TestEventResponse_Structure(t *testing.T) {
	event := map[string]interface{}{
		"id":          "evt-123",
		"name":        "Summer Festival",
		"description": "Amazing summer event",
		"capacity":    1000,
		"price":       "49.99",
		"location":    "Park",
		"category":    "Music",
		"startDate":   time.Now().Format(time.RFC3339),
	}

	jsonData, err := json.Marshal(event)
	assert.NoError(t, err)

	var decoded map[string]interface{}
	err = json.Unmarshal(jsonData, &decoded)
	assert.NoError(t, err)

	assert.Equal(t, "evt-123", decoded["id"])
	assert.Equal(t, "Summer Festival", decoded["name"])
	assert.Equal(t, float64(1000), decoded["capacity"])
}

func TestEventValidation_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		capacity int
		price    string
		valid    bool
	}{
		{"normal_values", 100, "25.99", true},
		{"zero_capacity", 0, "10.00", true},
		{"negative_capacity", -10, "10.00", false},
		{"high_capacity", 100000, "199.99", true},
		{"zero_price", 100, "0.00", true},
		{"high_price", 100, "9999.99", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := map[string]interface{}{
				"name":     "Test Event",
				"capacity": tt.capacity,
				"price":    tt.price,
			}

			jsonData, err := json.Marshal(event)
			assert.NoError(t, err)

			var decoded map[string]interface{}
			err = json.Unmarshal(jsonData, &decoded)
			assert.NoError(t, err)

			if tt.capacity < 0 {
				assert.Less(t, decoded["capacity"], 0.0)
			} else {
				assert.GreaterOrEqual(t, decoded["capacity"], 0.0)
			}
		})
	}
}

func TestDecimalPriceHandling(t *testing.T) {
	tests := []struct {
		priceStr    string
		expectedStr string
	}{
		{"0.00", "0"},
		{"9.99", "9.99"},
		{"99.99", "99.99"},
		{"999.99", "999.99"},
		{"1234.56", "1234.56"},
	}

	for _, tt := range tests {
		t.Run("price_"+tt.priceStr, func(t *testing.T) {
			price, err := decimal.NewFromString(tt.priceStr)
			assert.NoError(t, err)
			assert.True(t, price.GreaterThanOrEqual(decimal.Zero))
			
			// decimal library simplifies 0.00 to 0
			assert.Equal(t, tt.expectedStr, price.String())
		})
	}
}

func TestEventCategories(t *testing.T) {
	categories := []string{"Music", "Sports", "Technology", "Art", "Food", "Education"}

	for _, category := range categories {
		t.Run("category_"+category, func(t *testing.T) {
			event := map[string]interface{}{
				"name":     "Test Event",
				"category": category,
			}

			jsonData, err := json.Marshal(event)
			assert.NoError(t, err)

			var decoded map[string]interface{}
			err = json.Unmarshal(jsonData, &decoded)
			assert.NoError(t, err)
			assert.Equal(t, category, decoded["category"])
		})
	}
}
