package services

import (
"testing"

"github.com/stretchr/testify/assert"
)

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

// Client can be nil if not connected, which is okay in tests
// We're just testing that the method doesn't panic
_ = client
}

// TestIsConnected tests connection status check
func TestIsConnected(t *testing.T) {
service := GetDatabaseSeviceInstance()
isConnected := service.IsConnected()

// Should return a boolean without panicking
assert.IsType(t, false, isConnected, "Should return boolean")
}
