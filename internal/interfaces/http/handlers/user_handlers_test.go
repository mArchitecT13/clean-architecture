package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"clean-architecture/internal/domain/entities"
)

// MockUserUseCase is a mock implementation of UserUseCaseInterface
type MockUserUseCase struct {
	mock.Mock
}

func (m *MockUserUseCase) CreateUser(ctx context.Context, email, name string) (*entities.User, error) {
	args := m.Called(ctx, email, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserUseCase) GetUserByID(ctx context.Context, id string) (*entities.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserUseCase) UpdateUser(ctx context.Context, id, name, email string) (*entities.User, error) {
	args := m.Called(ctx, id, name, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserUseCase) DeleteUser(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserUseCase) ListUsers(ctx context.Context, limit, offset int) ([]*entities.User, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.User), args.Error(1)
}

func TestUserHandler_CreateUser(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    CreateUserRequest
		mockUser       *entities.User
		mockError      error
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "successful user creation",
			requestBody: CreateUserRequest{
				Email: "test@example.com",
				Name:  "Test User",
			},
			mockUser: &entities.User{
				ID:        "user_123",
				Email:     "test@example.com",
				Name:      "Test User",
				CreatedAt: entities.User{}.CreatedAt, // Will be set by entity
				UpdatedAt: entities.User{}.UpdatedAt, // Will be set by entity
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"status":  "success",
				"message": "User created successfully",
			},
		},
		{
			name: "invalid request body",
			requestBody: CreateUserRequest{
				Email: "test@example.com",
				Name:  "Test User",
			},
			mockUser:       nil,
			mockError:      assert.AnError,
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"status": "error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockUseCase := new(MockUserUseCase)
			handler := &UserHandler{
				userUseCase: mockUseCase,
			}

			// Mock expectations
			if tt.mockError == nil {
				mockUseCase.On("CreateUser", mock.Anything, tt.requestBody.Email, tt.requestBody.Name).
					Return(tt.mockUser, nil)
			} else {
				mockUseCase.On("CreateUser", mock.Anything, tt.requestBody.Email, tt.requestBody.Name).
					Return(nil, tt.mockError)
			}

			// Create request
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Execute
			handler.CreateUser(w, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			for key, expectedValue := range tt.expectedBody {
				assert.Equal(t, expectedValue, response[key])
			}

			mockUseCase.AssertExpectations(t)
		})
	}
}

func TestUserHandler_GetUser(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		mockUser       *entities.User
		mockError      error
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:   "successful user retrieval",
			userID: "user_123",
			mockUser: &entities.User{
				ID:        "user_123",
				Email:     "test@example.com",
				Name:      "Test User",
				CreatedAt: entities.User{}.CreatedAt,
				UpdatedAt: entities.User{}.UpdatedAt,
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"status": "success",
			},
		},
		{
			name:           "user not found",
			userID:         "user_123",
			mockUser:       nil,
			mockError:      assert.AnError,
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"status": "error",
			},
		},
		{
			name:           "missing user ID",
			userID:         "",
			mockUser:       nil,
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"status":  "error",
				"message": "User ID is required",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockUseCase := new(MockUserUseCase)
			handler := &UserHandler{
				userUseCase: mockUseCase,
			}

			// Mock expectations
			if tt.userID != "" {
				mockUseCase.On("GetUserByID", mock.Anything, tt.userID).
					Return(tt.mockUser, tt.mockError)
			}

			// Create request
			req := httptest.NewRequest("GET", "/users/"+tt.userID, nil)
			w := httptest.NewRecorder()

			// Setup Chi context for URL parameters
			rctx := chi.NewRouteContext()
			if tt.userID != "" {
				rctx.URLParams.Add("id", tt.userID)
			}
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			// Execute
			handler.GetUser(w, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			for key, expectedValue := range tt.expectedBody {
				assert.Equal(t, expectedValue, response[key])
			}

			mockUseCase.AssertExpectations(t)
		})
	}
}

func TestUserHandler_UpdateUser(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		requestBody    UpdateUserRequest
		mockUser       *entities.User
		mockError      error
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:   "successful user update",
			userID: "user_123",
			requestBody: UpdateUserRequest{
				Name:  "Updated User",
				Email: "updated@example.com",
			},
			mockUser: &entities.User{
				ID:        "user_123",
				Email:     "updated@example.com",
				Name:      "Updated User",
				CreatedAt: entities.User{}.CreatedAt,
				UpdatedAt: entities.User{}.UpdatedAt,
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"status":  "success",
				"message": "User updated successfully",
			},
		},
		{
			name:           "missing user ID",
			userID:         "",
			requestBody:    UpdateUserRequest{},
			mockUser:       nil,
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"status":  "error",
				"message": "User ID is required",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockUseCase := new(MockUserUseCase)
			handler := &UserHandler{
				userUseCase: mockUseCase,
			}

			// Mock expectations
			if tt.userID != "" {
				mockUseCase.On("UpdateUser", mock.Anything, tt.userID, tt.requestBody.Name, tt.requestBody.Email).
					Return(tt.mockUser, tt.mockError)
			}

			// Create request
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("PUT", "/users/"+tt.userID, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Setup Chi context for URL parameters
			rctx := chi.NewRouteContext()
			if tt.userID != "" {
				rctx.URLParams.Add("id", tt.userID)
			}
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			// Execute
			handler.UpdateUser(w, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			for key, expectedValue := range tt.expectedBody {
				assert.Equal(t, expectedValue, response[key])
			}

			mockUseCase.AssertExpectations(t)
		})
	}
}

func TestUserHandler_DeleteUser(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		mockError      error
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:           "successful user deletion",
			userID:         "user_123",
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"status":  "success",
				"message": "User deleted successfully",
			},
		},
		{
			name:           "user not found",
			userID:         "user_123",
			mockError:      assert.AnError,
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"status": "error",
			},
		},
		{
			name:           "missing user ID",
			userID:         "",
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"status":  "error",
				"message": "User ID is required",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockUseCase := new(MockUserUseCase)
			handler := &UserHandler{
				userUseCase: mockUseCase,
			}

			// Mock expectations
			if tt.userID != "" {
				mockUseCase.On("DeleteUser", mock.Anything, tt.userID).
					Return(tt.mockError)
			}

			// Create request
			req := httptest.NewRequest("DELETE", "/users/"+tt.userID, nil)
			w := httptest.NewRecorder()

			// Setup Chi context for URL parameters
			rctx := chi.NewRouteContext()
			if tt.userID != "" {
				rctx.URLParams.Add("id", tt.userID)
			}
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			// Execute
			handler.DeleteUser(w, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			for key, expectedValue := range tt.expectedBody {
				assert.Equal(t, expectedValue, response[key])
			}

			mockUseCase.AssertExpectations(t)
		})
	}
}

func TestUserHandler_ListUsers(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    string
		mockUsers      []*entities.User
		mockError      error
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:        "successful user listing",
			queryParams: "?limit=5&offset=0",
			mockUsers: []*entities.User{
				{
					ID:        "user_1",
					Email:     "user1@example.com",
					Name:      "User 1",
					CreatedAt: entities.User{}.CreatedAt,
					UpdatedAt: entities.User{}.UpdatedAt,
				},
				{
					ID:        "user_2",
					Email:     "user2@example.com",
					Name:      "User 2",
					CreatedAt: entities.User{}.CreatedAt,
					UpdatedAt: entities.User{}.UpdatedAt,
				},
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"status": "success",
			},
		},
		{
			name:           "database error",
			queryParams:    "?limit=5&offset=0",
			mockUsers:      nil,
			mockError:      assert.AnError,
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"status": "error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockUseCase := new(MockUserUseCase)
			handler := &UserHandler{
				userUseCase: mockUseCase,
			}

			// Mock expectations
			mockUseCase.On("ListUsers", mock.Anything, 5, 0).
				Return(tt.mockUsers, tt.mockError)

			// Create request
			req := httptest.NewRequest("GET", "/users"+tt.queryParams, nil)
			w := httptest.NewRecorder()

			// Execute
			handler.ListUsers(w, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			for key, expectedValue := range tt.expectedBody {
				assert.Equal(t, expectedValue, response[key])
			}

			mockUseCase.AssertExpectations(t)
		})
	}
}
