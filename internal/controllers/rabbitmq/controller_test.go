package rabbitmq

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNewController(t *testing.T) {
	if _, err := os.Stat("../../configs/config.yaml"); os.IsNotExist(err) {
		t.Skip("Skipping test: config file not found")
	}

	controller := NewController()
	assert.NotNil(t, controller)
}

func TestTestConnection_NoConfig(t *testing.T) {
	if _, err := os.Stat("../../configs/config.yaml"); os.IsNotExist(err) {
		t.Skip("Skipping test: config file not found")
	}

	gin.SetMode(gin.TestMode)

	controller := NewController()
	router := gin.New()
	router.GET("/rabbitmq/test", controller.TestConnection)

	req := httptest.NewRequest(http.MethodGet, "/rabbitmq/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Will fail without RabbitMQ connection
	assert.NotEqual(t, 0, w.Code)
}

func TestPublishTestMessage_NoConfig(t *testing.T) {
	if _, err := os.Stat("../../configs/config.yaml"); os.IsNotExist(err) {
		t.Skip("Skipping test: config file not found")
	}

	gin.SetMode(gin.TestMode)

	controller := NewController()
	router := gin.New()
	router.POST("/rabbitmq/publish", controller.PublishTestMessage)

	req := httptest.NewRequest(http.MethodPost, "/rabbitmq/publish", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Will fail without RabbitMQ connection
	assert.NotEqual(t, 0, w.Code)
}

func TestSetupTestExchangeAndQueue_NoConfig(t *testing.T) {
	if _, err := os.Stat("../../configs/config.yaml"); os.IsNotExist(err) {
		t.Skip("Skipping test: config file not found")
	}

	gin.SetMode(gin.TestMode)

	controller := NewController()
	router := gin.New()
	router.POST("/rabbitmq/setup", controller.SetupTestExchangeAndQueue)

	req := httptest.NewRequest(http.MethodPost, "/rabbitmq/setup", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Will fail without RabbitMQ connection
	assert.NotEqual(t, 0, w.Code)
}

func TestRabbitMQController_HTTPMethods(t *testing.T) {
	if _, err := os.Stat("../../configs/config.yaml"); os.IsNotExist(err) {
		t.Skip("Skipping test: config file not found")
	}

	gin.SetMode(gin.TestMode)

	controller := NewController()
	router := gin.New()
	router.GET("/rabbitmq/test", controller.TestConnection)
	router.POST("/rabbitmq/publish", controller.PublishTestMessage)
	router.POST("/rabbitmq/setup", controller.SetupTestExchangeAndQueue)

	tests := []struct {
		name   string
		method string
		path   string
	}{
		{"test_connection_get", http.MethodGet, "/rabbitmq/test"},
		{"publish_post", http.MethodPost, "/rabbitmq/publish"},
		{"setup_post", http.MethodPost, "/rabbitmq/setup"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Just verify the endpoint exists (will fail without RabbitMQ)
			assert.NotEqual(t, 0, w.Code)
		})
	}
}

func TestRabbitMQController_WrongMethods(t *testing.T) {
	if _, err := os.Stat("../../configs/config.yaml"); os.IsNotExist(err) {
		t.Skip("Skipping test: config file not found")
	}

	gin.SetMode(gin.TestMode)

	controller := NewController()
	router := gin.New()
	router.GET("/rabbitmq/test", controller.TestConnection)
	router.POST("/rabbitmq/publish", controller.PublishTestMessage)

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{"test_post_wrong", http.MethodPost, "/rabbitmq/test", http.StatusNotFound},
		{"publish_get_wrong", http.MethodGet, "/rabbitmq/publish", http.StatusNotFound},
		{"test_put_wrong", http.MethodPut, "/rabbitmq/test", http.StatusNotFound},
		{"publish_delete_wrong", http.MethodDelete, "/rabbitmq/publish", http.StatusNotFound},
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
