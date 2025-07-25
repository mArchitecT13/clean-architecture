package usecase

import (
	"context"
	"testing"

	"clean-architecture/internal/infrastructure/database"
	"clean-architecture/pkg/logger"
)

func TestUserUseCase_CreateUser(t *testing.T) {
	// Setup
	logger := logger.New()
	userRepo := database.NewMockUserRepository()
	userUseCase := NewUserUseCase(userRepo, logger)

	tests := []struct {
		name     string
		email    string
		userName string
		wantErr  bool
	}{
		{
			name:     "valid user creation",
			email:    "test@example.com",
			userName: "Test User",
			wantErr:  false,
		},
		{
			name:     "empty email",
			email:    "",
			userName: "Test User",
			wantErr:  true,
		},
		{
			name:     "empty name",
			email:    "test@example.com",
			userName: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := userUseCase.CreateUser(context.Background(), tt.email, tt.userName)

			if tt.wantErr {
				if err == nil {
					t.Errorf("CreateUser() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("CreateUser() unexpected error: %v", err)
				return
			}

			if user == nil {
				t.Errorf("CreateUser() expected user but got nil")
				return
			}

			if user.Email != tt.email {
				t.Errorf("CreateUser() email = %v, want %v", user.Email, tt.email)
			}

			if user.Name != tt.userName {
				t.Errorf("CreateUser() name = %v, want %v", user.Name, tt.userName)
			}
		})
	}
}

func TestUserUseCase_GetUserByID(t *testing.T) {
	// Setup
	logger := logger.New()
	userRepo := database.NewMockUserRepository()
	userUseCase := NewUserUseCase(userRepo, logger)

	// Create a test user first
	user, err := userUseCase.CreateUser(context.Background(), "test@example.com", "Test User")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	tests := []struct {
		name    string
		userID  string
		wantErr bool
	}{
		{
			name:    "existing user",
			userID:  user.ID,
			wantErr: false,
		},
		{
			name:    "non-existing user",
			userID:  "non-existing-id",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := userUseCase.GetUserByID(context.Background(), tt.userID)

			if tt.wantErr {
				if err == nil {
					t.Errorf("GetUserByID() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("GetUserByID() unexpected error: %v", err)
				return
			}

			if result == nil {
				t.Errorf("GetUserByID() expected user but got nil")
				return
			}

			if result.ID != tt.userID {
				t.Errorf("GetUserByID() user ID = %v, want %v", result.ID, tt.userID)
			}
		})
	}
}
