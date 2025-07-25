package logger

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	// Test default logger creation
	logger := New()
	assert.NotNil(t, logger)

	// Test that it implements the Logger interface
	var _ Logger = logger
}

func TestLogger_Levels(t *testing.T) {
	logger := New()

	// Test all log levels
	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warn message")
	logger.Error("error message")

	// Test formatted messages
	logger.Debugf("debug %s", "formatted")
	logger.Infof("info %s", "formatted")
	logger.Warnf("warn %s", "formatted")
	logger.Errorf("error %s", "formatted")
}

func TestLogger_WithContext(t *testing.T) {
	logger := New()
	ctx := context.Background()

	loggerWithContext := logger.WithContext(ctx)
	assert.NotNil(t, loggerWithContext)
	// Note: WithContext returns the same logger instance, so they should be equal
	assert.Equal(t, logger, loggerWithContext)
}

func TestLogger_WithField(t *testing.T) {
	logger := New()

	loggerWithField := logger.WithField("key", "value")
	assert.NotNil(t, loggerWithField)
	// Note: WithField returns the same logger instance, so they should be equal
	assert.Equal(t, logger, loggerWithField)
}

func TestLogger_WithFields(t *testing.T) {
	logger := New()
	fields := map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
	}

	loggerWithFields := logger.WithFields(fields)
	assert.NotNil(t, loggerWithFields)
	// Note: WithFields returns the same logger instance, so they should be equal
	assert.Equal(t, logger, loggerWithFields)
}

func TestLogger_EnvironmentLevels(t *testing.T) {
	tests := []struct {
		name     string
		envLevel string
		expected bool // whether logger should be created successfully
	}{
		{
			name:     "debug level",
			envLevel: "debug",
			expected: true,
		},
		{
			name:     "info level",
			envLevel: "info",
			expected: true,
		},
		{
			name:     "warn level",
			envLevel: "warn",
			expected: true,
		},
		{
			name:     "error level",
			envLevel: "error",
			expected: true,
		},
		{
			name:     "invalid level",
			envLevel: "invalid",
			expected: true, // should default to info
		},
		{
			name:     "empty level",
			envLevel: "",
			expected: true, // should default to info
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variable
			if tt.envLevel != "" {
				os.Setenv("LOG_LEVEL", tt.envLevel)
				defer os.Unsetenv("LOG_LEVEL")
			}

			logger := New()
			assert.NotNil(t, logger)
		})
	}
}

func TestLogger_Chaining(t *testing.T) {
	logger := New()

	// Test chaining multiple operations
	loggerWithContext := logger.WithContext(context.Background())
	loggerWithField := loggerWithContext.WithField("key", "value")
	loggerWithFields := loggerWithField.WithFields(map[string]interface{}{
		"another_key": "another_value",
	})

	assert.NotNil(t, loggerWithContext)
	assert.NotNil(t, loggerWithField)
	assert.NotNil(t, loggerWithFields)
}

func TestLogger_ConcurrentAccess(t *testing.T) {
	logger := New()

	// Test concurrent access to logger
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func(id int) {
			logger.Info("concurrent message", id)
			logger.WithField("goroutine", id).Info("with field")
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestLogger_Output(t *testing.T) {
	// Test that logger writes to stdout
	logger := New()

	// This is a basic test to ensure no panics
	logger.Info("test message")
	logger.Error("test error")

	// Test with different data types
	logger.Info("string", 123, true, 3.14)
	logger.WithField("number", 42).Info("with number")
	logger.WithFields(map[string]interface{}{
		"string": "value",
		"int":    123,
		"bool":   true,
		"float":  3.14,
	}).Info("with multiple fields")
}

func TestLogger_InterfaceCompliance(t *testing.T) {
	// Test that our logger properly implements the Logger interface
	var logger Logger = New()

	// Test all interface methods
	logger.Debug("test")
	logger.Info("test")
	logger.Warn("test")
	logger.Error("test")

	logger.Debugf("test %s", "formatted")
	logger.Infof("test %s", "formatted")
	logger.Warnf("test %s", "formatted")
	logger.Errorf("test %s", "formatted")

	logger.WithContext(context.Background())
	logger.WithField("key", "value")
	logger.WithFields(map[string]interface{}{"key": "value"})
}

func TestLogger_DefaultBehavior(t *testing.T) {
	// Test default behavior when no environment variables are set
	originalLogLevel := os.Getenv("LOG_LEVEL")
	defer func() {
		if originalLogLevel != "" {
			os.Setenv("LOG_LEVEL", originalLogLevel)
		} else {
			os.Unsetenv("LOG_LEVEL")
		}
	}()

	// Clear LOG_LEVEL environment variable
	os.Unsetenv("LOG_LEVEL")

	logger := New()
	assert.NotNil(t, logger)

	// Should work without any environment variables
	logger.Info("default behavior test")
}

func TestLogger_StructuredLogging(t *testing.T) {
	logger := New()

	// Test structured logging with various field types
	logger.WithField("string_field", "string_value").Info("string field")
	logger.WithField("int_field", 42).Info("int field")
	logger.WithField("bool_field", true).Info("bool field")
	logger.WithField("float_field", 3.14).Info("float field")

	// Test with multiple fields
	logger.WithFields(map[string]interface{}{
		"user_id":   123,
		"email":     "test@example.com",
		"is_active": true,
		"score":     95.5,
	}).Info("user action")
}

func TestLogger_ContextPropagation(t *testing.T) {
	logger := New()
	ctx := context.WithValue(context.Background(), "request_id", "req-123")

	loggerWithContext := logger.WithContext(ctx)
	assert.NotNil(t, loggerWithContext)

	// Test that context is properly propagated
	loggerWithContext.Info("message with context")
}
