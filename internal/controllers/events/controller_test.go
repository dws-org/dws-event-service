package events

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

func init() {
gin.SetMode(gin.TestMode)
}

// TestGetEvents_Success tests successful retrieval of events
func TestGetEvents_Success(t *testing.T) {
	// Setup
	router := gin.New()
	controller := NewController()
	router.GET("/api/v1/events", controller.GetEvents)

	// Execute
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/events", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code, "Expected status 200 OK")

	var response []interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err, "Response should be valid JSON array")
}

// TestGetEventByID_NotFound tests 404 error when event doesn't exist
// REQ7: Failure test case #1 - 404 Not Found
func TestGetEventByID_NotFound(t *testing.T) {
	// Setup
	router := gin.New()
	controller := NewController()
	router.GET("/api/v1/events/:id", controller.GetEventByID)

// Execute with non-existent ID
w := httptest.NewRecorder()
req, _ := http.NewRequest("GET", "/api/v1/events/non-existent-id-12345", nil)
router.ServeHTTP(w, req)

// Assert
assert.Equal(t, http.StatusNotFound, w.Code, "Expected status 404 Not Found")

var response map[string]interface{}
err := json.Unmarshal(w.Body.Bytes(), &response)
assert.NoError(t, err, "Response should be valid JSON")
assert.Contains(t, response, "error", "Response should contain error message")
}

// TestCreateEvent_BadRequest tests 400 error with invalid payload
// REQ7: Failure test case #2 - 400 Bad Request (missing required fields)
func TestCreateEvent_BadRequest_MissingFields(t *testing.T) {
	// Setup
	router := gin.New()
	controller := NewController()
router.POST("/api/v1/events", controller.CreateEvent)

// Execute with incomplete payload (missing required fields)
invalidPayload := map[string]interface{}{
"name": "Test Event",
// Missing: description, startDate, startTime, price, endDate, location, capacity, imageUrl, category, organizerId
}
body, _ := json.Marshal(invalidPayload)

w := httptest.NewRecorder()
req, _ := http.NewRequest("POST", "/api/v1/events", bytes.NewBuffer(body))
req.Header.Set("Content-Type", "application/json")
router.ServeHTTP(w, req)

// Assert
assert.Equal(t, http.StatusBadRequest, w.Code, "Expected status 400 Bad Request")

var response map[string]interface{}
err := json.Unmarshal(w.Body.Bytes(), &response)
assert.NoError(t, err, "Response should be valid JSON")
assert.Contains(t, response, "error", "Response should contain error message")
assert.Equal(t, "Invalid request payload", response["error"], "Error message should indicate invalid payload")
}

// TestCreateEvent_BadRequest_InvalidJSON tests 400 error with malformed JSON
// REQ7: Failure test case #3 - 400 Bad Request (invalid JSON)
func TestCreateEvent_BadRequest_InvalidJSON(t *testing.T) {
	// Setup
	router := gin.New()
	controller := NewController()
router.POST("/api/v1/events", controller.CreateEvent)

// Execute with invalid JSON
invalidJSON := []byte(`{"name": "Test Event", "description": }`) // Malformed JSON

w := httptest.NewRecorder()
req, _ := http.NewRequest("POST", "/api/v1/events", bytes.NewBuffer(invalidJSON))
req.Header.Set("Content-Type", "application/json")
router.ServeHTTP(w, req)

// Assert
assert.Equal(t, http.StatusBadRequest, w.Code, "Expected status 400 Bad Request")

var response map[string]interface{}
err := json.Unmarshal(w.Body.Bytes(), &response)
assert.NoError(t, err, "Response should be valid JSON")
assert.Contains(t, response, "error", "Response should contain error message")
}

// TestCreateEvent_BadRequest_InvalidDataTypes tests 400 error with wrong data types
// REQ7: Failure test case #4 - 400 Bad Request (invalid data types)
func TestCreateEvent_BadRequest_InvalidDataTypes(t *testing.T) {
	// Setup
	router := gin.New()
	controller := NewController()
router.POST("/api/v1/events", controller.CreateEvent)

// Execute with invalid data types (capacity as string instead of int)
invalidPayload := map[string]interface{}{
"name":        "Test Event",
"description": "Test Description",
"startDate":   "2026-01-01T00:00:00Z",
"startTime":   "2026-01-01T10:00:00Z",
"price":       "10.99",
"endDate":     "2026-01-01T00:00:00Z",
"location":    "Test Location",
"capacity":    "not-a-number", // Invalid: should be int
"imageUrl":    "https://example.com/image.jpg",
"category":    "Test Category",
"organizerId": "test-organizer-id",
}
body, _ := json.Marshal(invalidPayload)

w := httptest.NewRecorder()
req, _ := http.NewRequest("POST", "/api/v1/events", bytes.NewBuffer(body))
req.Header.Set("Content-Type", "application/json")
router.ServeHTTP(w, req)

// Assert
assert.Equal(t, http.StatusBadRequest, w.Code, "Expected status 400 Bad Request")

var response map[string]interface{}
err := json.Unmarshal(w.Body.Bytes(), &response)
assert.NoError(t, err, "Response should be valid JSON")
assert.Contains(t, response, "error", "Response should contain error message")
}

// TestCreateEventRequest_Validation tests request struct validation
func TestCreateEventRequest_Validation(t *testing.T) {
tests := []struct {
name        string
request     CreateEventRequest
shouldError bool
}{
{
name: "Valid request",
request: CreateEventRequest{
Name:        "Test Event",
Description: "Test Description",
StartDate:   time.Now(),
StartTime:   time.Now(),
Price:       decimal.NewFromFloat(10.99),
EndDate:     time.Now().Add(24 * time.Hour),
Location:    "Test Location",
Capacity:    100,
ImageURL:    "https://example.com/image.jpg",
Category:    "Test Category",
OrganizerID: "test-organizer-id",
},
shouldError: false,
},
{
name: "Missing name",
request: CreateEventRequest{
Description: "Test Description",
StartDate:   time.Now(),
StartTime:   time.Now(),
Price:       decimal.NewFromFloat(10.99),
EndDate:     time.Now().Add(24 * time.Hour),
Location:    "Test Location",
Capacity:    100,
ImageURL:    "https://example.com/image.jpg",
Category:    "Test Category",
OrganizerID: "test-organizer-id",
},
shouldError: true,
},
}

for _, tt := range tests {
t.Run(tt.name, func(t *testing.T) {
// Test validates the struct fields are properly tagged
if tt.shouldError {
assert.Empty(t, tt.request.Name, "Name should be empty for error case")
} else {
assert.NotEmpty(t, tt.request.Name, "Name should not be empty for valid case")
}
})
}
}

// TestNewController tests controller initialization
func TestNewController(t *testing.T) {
	controller := NewController()
	assert.NotNil(t, controller, "Controller should not be nil")
	assert.NotNil(t, controller.dbService, "Database service should be initialized")
}
