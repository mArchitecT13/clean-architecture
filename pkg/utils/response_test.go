package utils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSuccessResponse(t *testing.T) {
	data := map[string]interface{}{
		"id":   123,
		"name": "test",
	}
	message := "success message"

	response := SuccessResponse(data, message)

	assert.Equal(t, "success", response.Status)
	assert.Equal(t, message, response.Message)
	assert.Equal(t, data, response.Data)
	assert.False(t, response.Timestamp.IsZero())
}

func TestErrorResponse(t *testing.T) {
	message := "error message"

	response := ErrorResponse(message)

	assert.Equal(t, "error", response.Status)
	assert.Equal(t, message, response.Message)
	assert.Nil(t, response.Data)
	assert.False(t, response.Timestamp.IsZero())
}

func TestWriteJSON(t *testing.T) {
	w := httptest.NewRecorder()
	data := map[string]interface{}{
		"key": "value",
	}

	WriteJSON(w, http.StatusCreated, data)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "value", response["key"])
}

func TestWriteSuccess(t *testing.T) {
	w := httptest.NewRecorder()
	data := map[string]interface{}{
		"id":   123,
		"name": "test",
	}
	message := "success message"

	WriteSuccess(w, data, message)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response.Status)
	assert.Equal(t, message, response.Message)
	// JSON numbers are unmarshaled as float64, so we need to check the type
	assert.IsType(t, map[string]interface{}{}, response.Data)
	dataMap := response.Data.(map[string]interface{})
	assert.Equal(t, float64(123), dataMap["id"])
	assert.Equal(t, "test", dataMap["name"])
	assert.False(t, response.Timestamp.IsZero())
}

func TestWriteError(t *testing.T) {
	w := httptest.NewRecorder()
	message := "error message"

	WriteError(w, http.StatusBadRequest, message)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "error", response.Status)
	assert.Equal(t, message, response.Message)
	assert.Nil(t, response.Data)
	assert.False(t, response.Timestamp.IsZero())
}

func TestAPIResponse_JSONSerialization(t *testing.T) {
	tests := []struct {
		name     string
		response APIResponse
		expected map[string]interface{}
	}{
		{
			name: "success response with data",
			response: APIResponse{
				Status:    "success",
				Message:   "test message",
				Data:      map[string]interface{}{"key": "value"},
				Timestamp: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			expected: map[string]interface{}{
				"status":    "success",
				"message":   "test message",
				"data":      map[string]interface{}{"key": "value"},
				"timestamp": "2023-01-01T00:00:00Z",
			},
		},
		{
			name: "error response without data",
			response: APIResponse{
				Status:    "error",
				Message:   "error message",
				Data:      nil,
				Timestamp: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			expected: map[string]interface{}{
				"status":    "error",
				"message":   "error message",
				"timestamp": "2023-01-01T00:00:00Z",
			},
		},
		{
			name: "response without message",
			response: APIResponse{
				Status:    "success",
				Message:   "",
				Data:      map[string]interface{}{"id": 123},
				Timestamp: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			expected: map[string]interface{}{
				"status":    "success",
				"data":      map[string]interface{}{"id": float64(123)}, // JSON numbers are float64
				"timestamp": "2023-01-01T00:00:00Z",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Serialize to JSON
			jsonData, err := json.Marshal(tt.response)
			assert.NoError(t, err)

			// Deserialize back to map
			var result map[string]interface{}
			err = json.Unmarshal(jsonData, &result)
			assert.NoError(t, err)

			// Check expected fields
			for key, expectedValue := range tt.expected {
				assert.Equal(t, expectedValue, result[key])
			}

			// Check that optional fields are not present when empty
			if tt.response.Message == "" {
				_, exists := result["message"]
				assert.False(t, exists)
			}
			if tt.response.Data == nil {
				_, exists := result["data"]
				assert.False(t, exists)
			}
		})
	}
}

func TestWriteJSON_WithDifferentDataTypes(t *testing.T) {
	tests := []struct {
		name     string
		data     interface{}
		expected interface{}
	}{
		{
			name:     "string data",
			data:     "test string",
			expected: "test string",
		},
		{
			name:     "integer data",
			data:     123,
			expected: float64(123), // JSON numbers are float64
		},
		{
			name:     "boolean data",
			data:     true,
			expected: true,
		},
		{
			name:     "array data",
			data:     []string{"a", "b", "c"},
			expected: []interface{}{"a", "b", "c"},
		},
		{
			name:     "map data",
			data:     map[string]int{"a": 1, "b": 2},
			expected: map[string]interface{}{"a": float64(1), "b": float64(2)},
		},
		{
			name:     "nil data",
			data:     nil,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			WriteJSON(w, http.StatusOK, tt.data)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

			// For WriteJSON, the data is written directly, not wrapped in a response structure
			var result interface{}
			err := json.Unmarshal(w.Body.Bytes(), &result)
			assert.NoError(t, err)

			// The data should be in the response
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestWriteError_WithDifferentStatusCodes(t *testing.T) {
	tests := []struct {
		name         string
		statusCode   int
		message      string
		expectedCode int
	}{
		{
			name:         "bad request",
			statusCode:   http.StatusBadRequest,
			message:      "bad request",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "not found",
			statusCode:   http.StatusNotFound,
			message:      "not found",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "internal server error",
			statusCode:   http.StatusInternalServerError,
			message:      "internal error",
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "unauthorized",
			statusCode:   http.StatusUnauthorized,
			message:      "unauthorized",
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			WriteError(w, tt.statusCode, tt.message)

			assert.Equal(t, tt.expectedCode, w.Code)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

			var response APIResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "error", response.Status)
			assert.Equal(t, tt.message, response.Message)
		})
	}
}

func TestResponseHelpers_EdgeCases(t *testing.T) {
	// Test with empty message
	w := httptest.NewRecorder()
	WriteSuccess(w, nil, "")

	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response.Status)
	assert.Equal(t, "", response.Message)
	assert.Nil(t, response.Data)

	// Test with nil data
	w = httptest.NewRecorder()
	WriteSuccess(w, nil, "test")

	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response.Status)
	assert.Equal(t, "test", response.Message)
	assert.Nil(t, response.Data)

	// Test with complex data structure
	w = httptest.NewRecorder()
	complexData := map[string]interface{}{
		"user": map[string]interface{}{
			"id":     123,
			"name":   "test",
			"email":  "test@example.com",
			"active": true,
		},
		"metadata": []string{"tag1", "tag2"},
	}
	WriteSuccess(w, complexData, "complex data")

	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response.Status)
	assert.Equal(t, "complex data", response.Message)

	// Check that the data structure is preserved with correct types
	dataMap := response.Data.(map[string]interface{})
	userMap := dataMap["user"].(map[string]interface{})
	assert.Equal(t, float64(123), userMap["id"]) // JSON numbers are float64
	assert.Equal(t, "test", userMap["name"])
	assert.Equal(t, "test@example.com", userMap["email"])
	assert.Equal(t, true, userMap["active"])

	metadata := dataMap["metadata"].([]interface{})
	assert.Equal(t, "tag1", metadata[0])
	assert.Equal(t, "tag2", metadata[1])
}
