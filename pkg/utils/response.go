package utils

import (
	"encoding/json"
	"net/http"
	"time"
)

// APIResponse represents a standard API response
type APIResponse struct {
	Status    string      `json:"status"`
	Message   string      `json:"message,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// SuccessResponse creates a success response
func SuccessResponse(data interface{}, message string) APIResponse {
	return APIResponse{
		Status:    "success",
		Message:   message,
		Data:      data,
		Timestamp: time.Now(),
	}
}

// ErrorResponse creates an error response
func ErrorResponse(message string) APIResponse {
	return APIResponse{
		Status:    "error",
		Message:   message,
		Timestamp: time.Now(),
	}
}

// WriteJSON writes a JSON response
func WriteJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// WriteSuccess writes a success JSON response
func WriteSuccess(w http.ResponseWriter, data interface{}, message string) {
	response := SuccessResponse(data, message)
	WriteJSON(w, http.StatusOK, response)
}

// WriteError writes an error JSON response
func WriteError(w http.ResponseWriter, statusCode int, message string) {
	response := ErrorResponse(message)
	WriteJSON(w, statusCode, response)
}
