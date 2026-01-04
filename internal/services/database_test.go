package services

import (
"testing"

"github.com/stretchr/testify/assert"
)

// TestGetDatabaseServiceInstance tests singleton pattern
func TestGetDatabaseServiceInstance(t *testing.T) {
	// Skip if no real DATABASE_URL - requires PostgreSQL connection
	t.Skip("Skipping database test - requires actual PostgreSQL connection")
	
	instance1 := GetDatabaseSeviceInstance()
	instance2 := GetDatabaseSeviceInstance()

	assert.NotNil(t, instance1, "First instance should not be nil")
	assert.NotNil(t, instance2, "Second instance should not be nil")
	assert.Equal(t, instance1, instance2, "Both instances should be the same (singleton)")
}

// TestGetClient tests database client retrieval
func TestGetClient(t *testing.T) {
	// Skip if no real DATABASE_URL - requires PostgreSQL connection
	t.Skip("Skipping database test - requires actual PostgreSQL connection")
	
	service := GetDatabaseSeviceInstance()
	client := service.GetClient()

	// Client should not be nil once service is instantiated
	assert.NotNil(t, client, "Client should be initialized")
}
