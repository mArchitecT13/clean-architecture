package database

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"clean-architecture/internal/domain/entities"
	"clean-architecture/internal/domain/repositories"
)

// MockUserRepository implements UserRepository interface for testing
type MockUserRepository struct {
	users map[string]*entities.User
	mutex sync.RWMutex
}

// NewMockUserRepository creates a new mock user repository
func NewMockUserRepository() repositories.UserRepository {
	return &MockUserRepository{
		users: make(map[string]*entities.User),
	}
}

// Create creates a new user
func (r *MockUserRepository) Create(ctx context.Context, user *entities.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Generate ID if not set
	if user.ID == "" {
		user.ID = fmt.Sprintf("user_%d", time.Now().UnixNano())
	}

	// Set timestamps if not set
	now := time.Now()
	if user.CreatedAt.IsZero() {
		user.CreatedAt = now
	}
	if user.UpdatedAt.IsZero() {
		user.UpdatedAt = now
	}

	// Check if user with same email exists
	for _, existingUser := range r.users {
		if existingUser.Email == user.Email {
			return errors.New("user with this email already exists")
		}
	}

	// Store user
	r.users[user.ID] = user
	return nil
}

// GetByID retrieves a user by ID
func (r *MockUserRepository) GetByID(ctx context.Context, id string) (*entities.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, nil
	}

	// Return a copy to avoid external modifications
	return &entities.User{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

// GetByEmail retrieves a user by email
func (r *MockUserRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, user := range r.users {
		if user.Email == email {
			// Return a copy to avoid external modifications
			return &entities.User{
				ID:        user.ID,
				Email:     user.Email,
				Name:      user.Name,
				CreatedAt: user.CreatedAt,
				UpdatedAt: user.UpdatedAt,
			}, nil
		}
	}

	return nil, nil
}

// Update updates a user
func (r *MockUserRepository) Update(ctx context.Context, user *entities.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	existingUser, exists := r.users[user.ID]
	if !exists {
		return errors.New("user not found")
	}

	// Update the user with current timestamp
	now := time.Now()
	r.users[user.ID] = &entities.User{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: existingUser.CreatedAt,
		UpdatedAt: now,
	}

	return nil
}

// Delete deletes a user
func (r *MockUserRepository) Delete(ctx context.Context, id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.users[id]; !exists {
		return errors.New("user not found")
	}

	delete(r.users, id)
	return nil
}

// List retrieves a list of users
func (r *MockUserRepository) List(ctx context.Context, limit, offset int) ([]*entities.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var users []*entities.User
	count := 0

	for _, user := range r.users {
		if count >= offset {
			if len(users) >= limit {
				break
			}
			// Return a copy to avoid external modifications
			users = append(users, &entities.User{
				ID:        user.ID,
				Email:     user.Email,
				Name:      user.Name,
				CreatedAt: user.CreatedAt,
				UpdatedAt: user.UpdatedAt,
			})
		}
		count++
	}

	return users, nil
}
