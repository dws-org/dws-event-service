package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLogrusLogger(t *testing.T) {
	logger := NewLogrusLogger()
	
	assert.NotNil(t, logger)
	assert.NotNil(t, logger.Out, "Logger output should be set")
}

func TestLoggerWriter(t *testing.T) {
	logger := NewLogrusLogger()
	writer := logger.Writer()
	
	assert.NotNil(t, writer, "Writer should not be nil")
	
	// Test writing to logger
	n, err := writer.Write([]byte("test log message"))
	assert.NoError(t, err)
	assert.Greater(t, n, 0, "Should write bytes")
}

func TestLoggerLevels(t *testing.T) {
	logger := NewLogrusLogger()
	
	// Test different log levels don't panic
	assert.NotPanics(t, func() {
		logger.Trace("Trace message")
		logger.Debug("Debug message")
		logger.Info("Info message")
		logger.Warn("Warning message")
		logger.Error("Error message")
	})
}

func TestLoggerFormatting(t *testing.T) {
	logger := NewLogrusLogger()
	
	// Test formatted logging
	assert.NotPanics(t, func() {
		logger.Infof("Formatted message: %s %d", "test", 123)
		logger.Warnf("Warning with value: %v", map[string]int{"count": 5})
		logger.Errorf("Error code: %d", 500)
	})
}

func TestLoggerWithFields(t *testing.T) {
	logger := NewLogrusLogger()
	
	// Test structured logging with fields
	assert.NotPanics(t, func() {
		logger.WithField("user_id", "123").Info("User logged in")
		logger.WithFields(map[string]interface{}{
			"method": "POST",
			"path":   "/api/events",
			"status": 201,
		}).Info("HTTP request")
	})
}

func TestMultipleLoggerInstances(t *testing.T) {
	logger1 := NewLogrusLogger()
	logger2 := NewLogrusLogger()
	
	assert.NotNil(t, logger1)
	assert.NotNil(t, logger2)
	
	// Each call creates a new instance
	// They may not be equal objects
	assert.NotNil(t, logger1.Out)
	assert.NotNil(t, logger2.Out)
}
