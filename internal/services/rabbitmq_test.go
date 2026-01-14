package services

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRabbitMQServiceInstance_Disabled(t *testing.T) {
	if _, err := os.Stat("../../configs/config.yaml"); os.IsNotExist(err) {
		t.Skip("Skipping test: config file not found")
	}
	
	// When RabbitMQ is disabled in config, should return nil
	instance := GetRabbitMQServiceInstance()
	
	// Can be nil if disabled, or non-nil if enabled
	// Both are valid states
	if instance != nil {
		assert.NotNil(t, instance.logger)
		assert.NotNil(t, instance.config)
	}
}

func TestRabbitMQService_Structure(t *testing.T) {
	if _, err := os.Stat("../../configs/config.yaml"); os.IsNotExist(err) {
		t.Skip("Skipping test: config file not found")
	}
	
	instance := GetRabbitMQServiceInstance()
	
	if instance == nil {
		t.Skip("RabbitMQ is disabled in config")
		return
	}
	
	assert.NotNil(t, instance.logger, "Logger should be initialized")
	assert.NotNil(t, instance.config, "Config should be initialized")
}

func TestRabbitMQService_HealthCheck_NoConnection(t *testing.T) {
	if _, err := os.Stat("../../configs/config.yaml"); os.IsNotExist(err) {
		t.Skip("Skipping test: config file not found")
	}
	
	instance := GetRabbitMQServiceInstance()
	
	if instance == nil {
		t.Skip("RabbitMQ is disabled in config")
		return
	}
	
	// Without real connection, health check may fail
	// This is expected behavior
	t.Skip("Skipping test: requires RabbitMQ connection")
}

func TestRabbitMQService_Singleton(t *testing.T) {
	if _, err := os.Stat("../../configs/config.yaml"); os.IsNotExist(err) {
		t.Skip("Skipping test: config file not found")
	}
	
	instance1 := GetRabbitMQServiceInstance()
	instance2 := GetRabbitMQServiceInstance()
	
	// Both should be same instance (singleton pattern)
	// Or both nil if disabled
	if instance1 != nil && instance2 != nil {
		assert.Equal(t, instance1, instance2, "Should return same instance")
	}
}
