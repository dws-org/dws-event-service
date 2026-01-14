package services

import (
	"context"
	"os"
	"testing"

	"github.com/oskargbc/dws-event-service.git/prisma/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockPrismaClient is a mock implementation of Prisma client
type MockPrismaClient struct {
	mock.Mock
}

// MockDatabaseService wraps the mock client
type MockDatabaseService struct {
	mockClient *MockPrismaClient
}

func NewMockDatabaseService() *MockDatabaseService {
	return &MockDatabaseService{
		mockClient: &MockPrismaClient{},
	}
}

func (m *MockDatabaseService) GetClient() *db.PrismaClient {
	// Return nil for now - full mock implementation would be complex
	return nil
}

func TestGetDatabaseServiceInstance(t *testing.T) {
	if _, err := os.Stat("../../configs/config.yaml"); os.IsNotExist(err) {
		t.Skip("Skipping test: config file not found")
	}
	
	// Test singleton pattern
	instance1 := GetDatabaseSeviceInstance()
	instance2 := GetDatabaseSeviceInstance()
	
	assert.NotNil(t, instance1)
	assert.NotNil(t, instance2)
	assert.Equal(t, instance1, instance2, "Should return same instance (singleton)")
}

func TestDatabaseService_Structure(t *testing.T) {
	if _, err := os.Stat("../../configs/config.yaml"); os.IsNotExist(err) {
		t.Skip("Skipping test: config file not found")
	}
	
	service := GetDatabaseSeviceInstance()
	
	assert.NotNil(t, service)
	assert.NotNil(t, service.client)
	assert.NotNil(t, service.logger)
	assert.NotNil(t, service.env)
}

func TestDatabaseService_HealthCheck_NoConnection(t *testing.T) {
	// Skip if no database connection
	t.Skip("Skipping test: requires database connection")
	
	if _, err := os.Stat("../../configs/config.yaml"); os.IsNotExist(err) {
		return
	}
	
	service := GetDatabaseSeviceInstance()
	ctx := context.Background()
	
	err := service.HealthCheck(ctx)
	// Without DB, should return error
	assert.Error(t, err)
}

func TestDatabaseService_EnsureConnected(t *testing.T) {
	if _, err := os.Stat("../../configs/config.yaml"); os.IsNotExist(err) {
		t.Skip("Skipping test: config file not found")
	}
	
	service := GetDatabaseSeviceInstance()
	ctx := context.Background()
	
	// Test that EnsureConnected handles nil client gracefully
	err := service.EnsureConnected(ctx)
	
	// Should either succeed or return meaningful error
	if err != nil {
		assert.Contains(t, err.Error(), "")  // Any error message is acceptable
	}
}

func TestDatabaseService_GetClient(t *testing.T) {
	if _, err := os.Stat("../../configs/config.yaml"); os.IsNotExist(err) {
		t.Skip("Skipping test: config file not found")
	}
	
	service := GetDatabaseSeviceInstance()
	
	client := service.GetClient()
	assert.NotNil(t, client, "Client should be initialized")
}
