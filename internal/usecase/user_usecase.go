package usecase

import (
	"context"
	"errors"
	"fmt"

	"clean-architecture/internal/domain/entities"
	"clean-architecture/internal/domain/repositories"
	"clean-architecture/pkg/logger"
)

// UserUseCase implements business logic for user operations
type UserUseCase struct {
	userRepo repositories.UserRepository
	logger   logger.Logger
}

// NewUserUseCase creates a new user use case instance
func NewUserUseCase(userRepo repositories.UserRepository, logger logger.Logger) *UserUseCase {
	return &UserUseCase{
		userRepo: userRepo,
		logger:   logger,
	}
}

// CreateUser creates a new user
func (uc *UserUseCase) CreateUser(ctx context.Context, email, name string) (*entities.User, error) {
	uc.logger.WithField("email", email).Info("Creating new user")

	// Validate input
	if email == "" {
		return nil, errors.New("email is required")
	}
	if name == "" {
		return nil, errors.New("name is required")
	}

	// Check if user already exists
	existingUser, err := uc.userRepo.GetByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Create new user
	user := entities.NewUser(email, name)

	// Save to repository
	err = uc.userRepo.Create(ctx, user)
	if err != nil {
		uc.logger.WithField("error", err.Error()).Error("Failed to create user")
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	uc.logger.WithField("user_id", user.ID).Info("User created successfully")
	return user, nil
}

// GetUserByID retrieves a user by ID
func (uc *UserUseCase) GetUserByID(ctx context.Context, id string) (*entities.User, error) {
	uc.logger.WithField("user_id", id).Debug("Getting user by ID")

	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		uc.logger.WithField("error", err.Error()).Error("Failed to get user by ID")
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

// UpdateUser updates user information
func (uc *UserUseCase) UpdateUser(ctx context.Context, id, name, email string) (*entities.User, error) {
	uc.logger.WithField("user_id", id).Info("Updating user")

	// Get existing user
	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		uc.logger.WithField("error", err.Error()).Error("Failed to get user for update")
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	// Update fields if provided
	if name != "" {
		user.UpdateName(name)
	}
	if email != "" {
		user.UpdateEmail(email)
	}

	// Save changes
	err = uc.userRepo.Update(ctx, user)
	if err != nil {
		uc.logger.WithField("error", err.Error()).Error("Failed to update user")
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	uc.logger.WithField("user_id", user.ID).Info("User updated successfully")
	return user, nil
}

// DeleteUser deletes a user
func (uc *UserUseCase) DeleteUser(ctx context.Context, id string) error {
	uc.logger.WithField("user_id", id).Info("Deleting user")

	err := uc.userRepo.Delete(ctx, id)
	if err != nil {
		uc.logger.WithField("error", err.Error()).Error("Failed to delete user")
		return fmt.Errorf("failed to delete user: %w", err)
	}

	uc.logger.WithField("user_id", id).Info("User deleted successfully")
	return nil
}

// ListUsers retrieves a list of users
func (uc *UserUseCase) ListUsers(ctx context.Context, limit, offset int) ([]*entities.User, error) {
	uc.logger.WithFields(map[string]interface{}{
		"limit":  limit,
		"offset": offset,
	}).Debug("Listing users")

	users, err := uc.userRepo.List(ctx, limit, offset)
	if err != nil {
		uc.logger.WithField("error", err.Error()).Error("Failed to list users")
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	return users, nil
}
