package handlers

import (
	"net/http"
	"time"

	"github.com/go-chi/render"
)

// Response represents a standard API response
type Response struct {
	Status    string      `json:"status"`
	Message   string      `json:"message,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// HealthCheck handles health check requests
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	response := Response{
		Status:    "success",
		Message:   "Service is healthy",
		Timestamp: time.Now(),
	}

	render.JSON(w, r, response)
}

// RootHandler handles root API requests
func RootHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{
		Status:  "success",
		Message: "Clean Architecture API",
		Data: map[string]interface{}{
			"version": "1.0.0",
			"docs":    "/docs",
		},
		Timestamp: time.Now(),
	}

	render.JSON(w, r, response)
}

// NotFoundHandler handles 404 requests
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{
		Status:    "error",
		Message:   "Endpoint not found",
		Timestamp: time.Now(),
	}

	w.WriteHeader(http.StatusNotFound)
	render.JSON(w, r, response)
}

// MethodNotAllowedHandler handles 405 requests
func MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{
		Status:    "error",
		Message:   "Method not allowed",
		Timestamp: time.Now(),
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
	render.JSON(w, r, response)
}

// UserResponse represents a user response for Swagger
// swagger:response UserResponse
type UserResponse struct {
	// in: body
	Body struct {
		Status    string      `json:"status"`
		Message   string      `json:"message,omitempty"`
		Data      interface{} `json:"data,omitempty"`
		Timestamp string      `json:"timestamp"`
	}
}

// ErrorResponse represents an error response for Swagger
// swagger:response ErrorResponse
type ErrorResponse struct {
	// in: body
	Body struct {
		Status    string `json:"status"`
		Message   string `json:"message"`
		Timestamp string `json:"timestamp"`
	}
}

// SuccessResponse represents a success response for Swagger
// swagger:response SuccessResponse
type SuccessResponse struct {
	// in: body
	Body struct {
		Status    string `json:"status"`
		Message   string `json:"message"`
		Timestamp string `json:"timestamp"`
	}
}
