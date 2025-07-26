package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"clean-architecture/internal/usecase"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	userUseCase usecase.UserUseCaseInterface
}

// NewUserHandler creates a new user handler
func NewUserHandler(userUseCase usecase.UserUseCaseInterface) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
	}
}

// CreateUserRequest represents the request body for creating a user
type CreateUserRequest struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

// UpdateUserRequest represents the request body for updating a user
type UpdateUserRequest struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}

// CreateUser godoc
// @Summary      Create a new user
// @Description  Create a new user with email and name
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      CreateUserRequest  true  "User info"
// @Success      200   {object}  UserResponse
// @Failure      400   {object}  ErrorResponse
// @Router       /api/v1/users [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.JSON(w, r, Response{
			Status:    "error",
			Message:   "Invalid request body",
			Timestamp: time.Now(),
		})
		return
	}

	user, err := h.userUseCase.CreateUser(r.Context(), req.Email, req.Name)
	if err != nil {
		render.JSON(w, r, Response{
			Status:    "error",
			Message:   err.Error(),
			Timestamp: time.Now(),
		})
		return
	}

	render.JSON(w, r, Response{
		Status:    "success",
		Message:   "User created successfully",
		Data:      user,
		Timestamp: time.Now(),
	})
}

// GetUser godoc
// @Summary      Get user by ID
// @Description  Get a user by their ID
// @Tags         users
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  UserResponse
// @Failure      404  {object}  ErrorResponse
// @Router       /api/v1/users/{id} [get]
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	if userID == "" {
		render.JSON(w, r, Response{
			Status:    "error",
			Message:   "User ID is required",
			Timestamp: time.Now(),
		})
		return
	}

	user, err := h.userUseCase.GetUserByID(r.Context(), userID)
	if err != nil {
		render.JSON(w, r, Response{
			Status:    "error",
			Message:   err.Error(),
			Timestamp: time.Now(),
		})
		return
	}

	render.JSON(w, r, Response{
		Status:    "success",
		Data:      user,
		Timestamp: time.Now(),
	})
}

// UpdateUser godoc
// @Summary      Update a user
// @Description  Update a user's information
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id    path      string             true  "User ID"
// @Param        user  body      UpdateUserRequest  true  "User info"
// @Success      200   {object}  UserResponse
// @Failure      400   {object}  ErrorResponse
// @Router       /api/v1/users/{id} [put]
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	if userID == "" {
		render.JSON(w, r, Response{
			Status:    "error",
			Message:   "User ID is required",
			Timestamp: time.Now(),
		})
		return
	}

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.JSON(w, r, Response{
			Status:    "error",
			Message:   "Invalid request body",
			Timestamp: time.Now(),
		})
		return
	}

	user, err := h.userUseCase.UpdateUser(r.Context(), userID, req.Name, req.Email)
	if err != nil {
		render.JSON(w, r, Response{
			Status:    "error",
			Message:   err.Error(),
			Timestamp: time.Now(),
		})
		return
	}

	render.JSON(w, r, Response{
		Status:    "success",
		Message:   "User updated successfully",
		Data:      user,
		Timestamp: time.Now(),
	})
}

// DeleteUser godoc
// @Summary      Delete a user
// @Description  Delete a user by their ID
// @Tags         users
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  SuccessResponse
// @Failure      404  {object}  ErrorResponse
// @Router       /api/v1/users/{id} [delete]
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	if userID == "" {
		render.JSON(w, r, Response{
			Status:    "error",
			Message:   "User ID is required",
			Timestamp: time.Now(),
		})
		return
	}

	err := h.userUseCase.DeleteUser(r.Context(), userID)
	if err != nil {
		render.JSON(w, r, Response{
			Status:    "error",
			Message:   err.Error(),
			Timestamp: time.Now(),
		})
		return
	}

	render.JSON(w, r, Response{
		Status:    "success",
		Message:   "User deleted successfully",
		Timestamp: time.Now(),
	})
}

// ListUsers godoc
// @Summary      List all users
// @Description  Get a list of all users
// @Tags         users
// @Produce      json
// @Success      200  {array}   UserResponse
// @Router       /api/v1/users [get]
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 10 // default limit
	offset := 0 // default offset

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	users, err := h.userUseCase.ListUsers(r.Context(), limit, offset)
	if err != nil {
		render.JSON(w, r, Response{
			Status:    "error",
			Message:   err.Error(),
			Timestamp: time.Now(),
		})
		return
	}

	render.JSON(w, r, Response{
		Status:    "success",
		Data:      users,
		Timestamp: time.Now(),
	})
}
