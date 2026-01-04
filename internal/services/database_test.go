package services

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// Run tests
	code := m.Run()
	
	// Clean up: disconnect database if instance exists
	if databaseInstance != nil {
		databaseInstance.DbDisconnect()
	}
	
	// Exit with test result code
	os.Exit(code)
}

// TestGetDatabaseServiceInstance tests singleton pattern
func TestGetDatabaseServiceInstance(t *testing.T) {
	instance1 := GetDatabaseSeviceInstance()
	instance2 := GetDatabaseSeviceInstance()

	assert.NotNil(t, instance1, "First instance should not be nil")
	assert.NotNil(t, instance2, "Second instance should not be nil")
	assert.Equal(t, instance1, instance2, "Both instances should be the same (singleton)")
}

// TestGetClient tests database client retrieval
func TestGetClient(t *testing.T) {
	service := GetDatabaseSeviceInstance()
	client := service.GetClient()

	// Client should not be nil once service is instantiated
	assert.NotNil(t, client, "Client should be initialized")
}
