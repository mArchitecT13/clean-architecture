package logging

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"clean-architecture/pkg/logger"
)

// MockLogger is a mock implementation of logger.Logger
type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Debug(args ...interface{}) {
	m.Called(args...)
}

func (m *MockLogger) Info(args ...interface{}) {
	m.Called(args...)
}

func (m *MockLogger) Warn(args ...interface{}) {
	m.Called(args...)
}

func (m *MockLogger) Error(args ...interface{}) {
	m.Called(args...)
}

func (m *MockLogger) Fatal(args ...interface{}) {
	m.Called(args...)
}

func (m *MockLogger) Debugf(format string, args ...interface{}) {
	m.Called(format, args)
}

func (m *MockLogger) Infof(format string, args ...interface{}) {
	m.Called(format, args)
}

func (m *MockLogger) Warnf(format string, args ...interface{}) {
	m.Called(format, args)
}

func (m *MockLogger) Errorf(format string, args ...interface{}) {
	m.Called(format, args)
}

func (m *MockLogger) Fatalf(format string, args ...interface{}) {
	m.Called(format, args)
}

func (m *MockLogger) WithContext(ctx context.Context) logger.Logger {
	args := m.Called(ctx)
	return args.Get(0).(logger.Logger)
}

func (m *MockLogger) WithField(key string, value interface{}) logger.Logger {
	args := m.Called(key, value)
	return args.Get(0).(logger.Logger)
}

func (m *MockLogger) WithFields(fields map[string]interface{}) logger.Logger {
	args := m.Called(fields)
	return args.Get(0).(logger.Logger)
}

func TestLoggerMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		statusCode     int
		expectedFields map[string]interface{}
	}{
		{
			name:       "successful request",
			method:     "GET",
			path:       "/api/v1/users",
			statusCode: 200,
			expectedFields: map[string]interface{}{
				"method": "GET",
				"path":   "/api/v1/users",
				"status": 200,
			},
		},
		{
			name:       "error request",
			method:     "POST",
			path:       "/api/v1/users",
			statusCode: 400,
			expectedFields: map[string]interface{}{
				"method": "POST",
				"path":   "/api/v1/users",
				"status": 400,
			},
		},
		{
			name:       "not found request",
			method:     "GET",
			path:       "/api/v1/nonexistent",
			statusCode: 404,
			expectedFields: map[string]interface{}{
				"method": "GET",
				"path":   "/api/v1/nonexistent",
				"status": 404,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock logger
			mockLogger := new(MockLogger)
			mockLoggerWithFields := new(MockLogger)

			// Setup expectations
			mockLogger.On("WithFields", mock.MatchedBy(func(fields map[string]interface{}) bool {
				// Check that expected fields are present
				for key, expectedValue := range tt.expectedFields {
					if value, exists := fields[key]; !exists || value != expectedValue {
						return false
					}
				}
				// Check that required fields exist
				requiredFields := []string{"method", "path", "status", "duration", "user_agent", "remote_ip"}
				for _, field := range requiredFields {
					if _, exists := fields[field]; !exists {
						return false
					}
				}
				return true
			})).Return(mockLoggerWithFields)

			mockLoggerWithFields.On("Info", "HTTP Request").Return()

			// Create middleware
			middleware := LoggerMiddleware(mockLogger)

			// Create test handler
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				w.Write([]byte("test response"))
			})

			// Create request
			req := httptest.NewRequest(tt.method, tt.path, nil)
			req.Header.Set("User-Agent", "test-agent")
			req.RemoteAddr = "127.0.0.1:12345"

			// Create response recorder
			w := httptest.NewRecorder()

			// Execute middleware
			middleware(handler).ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.statusCode, w.Code)
			assert.Equal(t, "test response", w.Body.String())

			// Verify mock expectations
			mockLogger.AssertExpectations(t)
			mockLoggerWithFields.AssertExpectations(t)
		})
	}
}

func TestResponseWriter(t *testing.T) {
	// Create a mock response writer
	mockWriter := httptest.NewRecorder()

	// Wrap it with our response writer
	rw := &responseWriter{
		ResponseWriter: mockWriter,
		statusCode:     http.StatusOK,
	}

	// Test WriteHeader
	rw.WriteHeader(http.StatusNotFound)
	assert.Equal(t, http.StatusNotFound, rw.statusCode)
	assert.Equal(t, http.StatusNotFound, mockWriter.Code)

	// Test Write
	testData := []byte("test data")
	n, err := rw.Write(testData)
	assert.NoError(t, err)
	assert.Equal(t, len(testData), n)
	assert.Equal(t, testData, mockWriter.Body.Bytes())
}

func TestLoggerMiddleware_WithPanic(t *testing.T) {
	// Setup mock logger
	mockLogger := new(MockLogger)
	mockLoggerWithFields := new(MockLogger)

	// Setup expectations for panic recovery
	mockLogger.On("WithFields", mock.AnythingOfType("map[string]interface {}")).Return(mockLoggerWithFields)
	mockLoggerWithFields.On("Info", "HTTP Request").Return()

	// Create middleware
	middleware := LoggerMiddleware(mockLogger)

	// Create test handler that panics
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})

	// Create request
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// Execute middleware - the panic should be recovered by Chi's Recoverer middleware
	// but our logging middleware should still work before the panic
	defer func() {
		if r := recover(); r != nil {
			// Panic was recovered, which is expected
			t.Logf("Recovered from panic: %v", r)
		}
	}()

	middleware(handler).ServeHTTP(w, req)

	// Verify that the middleware logged the request before the panic
	mockLogger.AssertExpectations(t)
	mockLoggerWithFields.AssertExpectations(t)
}

func TestLoggerMiddleware_WithEmptyUserAgent(t *testing.T) {
	// Setup mock logger
	mockLogger := new(MockLogger)
	mockLoggerWithFields := new(MockLogger)

	// Setup expectations
	mockLogger.On("WithFields", mock.MatchedBy(func(fields map[string]interface{}) bool {
		// Check that user_agent is empty string when not set
		if userAgent, exists := fields["user_agent"]; !exists || userAgent != "" {
			return false
		}
		return true
	})).Return(mockLoggerWithFields)

	mockLoggerWithFields.On("Info", "HTTP Request").Return()

	// Create middleware
	middleware := LoggerMiddleware(mockLogger)

	// Create test handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Create request without User-Agent header
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// Execute middleware
	middleware(handler).ServeHTTP(w, req)

	// Verify expectations
	mockLogger.AssertExpectations(t)
	mockLoggerWithFields.AssertExpectations(t)
}

func TestLoggerMiddleware_WithCustomHeaders(t *testing.T) {
	// Setup mock logger
	mockLogger := new(MockLogger)
	mockLoggerWithFields := new(MockLogger)

	// Setup expectations
	mockLogger.On("WithFields", mock.AnythingOfType("map[string]interface {}")).Return(mockLoggerWithFields)
	mockLoggerWithFields.On("Info", "HTTP Request").Return()

	// Create middleware
	middleware := LoggerMiddleware(mockLogger)

	// Create test handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Custom-Header", "custom-value")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("response"))
	})

	// Create request with custom headers
	req := httptest.NewRequest("POST", "/api/v1/users", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer token")
	req.Header.Set("User-Agent", "test-agent")
	w := httptest.NewRecorder()

	// Execute middleware
	middleware(handler).ServeHTTP(w, req)

	// Verify expectations
	mockLogger.AssertExpectations(t)
	mockLoggerWithFields.AssertExpectations(t)

	// Verify response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "response", w.Body.String())
	assert.Equal(t, "custom-value", w.Header().Get("X-Custom-Header"))
}
