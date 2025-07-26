package database

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"clean-architecture/internal/domain/entities"
	"clean-architecture/internal/domain/repositories"

	"gorm.io/gorm"
)

// PostgresUserRepository implements UserRepository interface using PostgreSQL
type PostgresUserRepository struct {
	db *gorm.DB
}

// NewPostgresUserRepository creates a new PostgreSQL user repository
func NewPostgresUserRepository(db *gorm.DB) repositories.UserRepository {
	return &PostgresUserRepository{db: db}
}

// Create creates a new user
func (r *PostgresUserRepository) Create(ctx context.Context, user *entities.User) error {
	// Check if user with same email exists
	var existingUser entities.User
	if err := r.db.WithContext(ctx).Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		return errors.New("user with this email already exists")
	}

	// Generate ID if not set
	if user.ID == "" {
		user.ID = generateID()
	}

	// Set timestamps if not set
	now := time.Now()
	if user.CreatedAt.IsZero() {
		user.CreatedAt = now
	}
	if user.UpdatedAt.IsZero() {
		user.UpdatedAt = now
	}

	return r.db.WithContext(ctx).Create(user).Error
}

// GetByID retrieves a user by ID
func (r *PostgresUserRepository) GetByID(ctx context.Context, id string) (*entities.User, error) {
	var user entities.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	var user entities.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// Update updates a user
func (r *PostgresUserRepository) Update(ctx context.Context, user *entities.User) error {
	// Check if user exists
	var existingUser entities.User
	if err := r.db.WithContext(ctx).Where("id = ?", user.ID).First(&existingUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	// Update the user with current timestamp
	user.UpdatedAt = time.Now()
	user.CreatedAt = existingUser.CreatedAt // Preserve original creation time

	return r.db.WithContext(ctx).Save(user).Error
}

// Delete deletes a user
func (r *PostgresUserRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&entities.User{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}

// List retrieves a list of users
func (r *PostgresUserRepository) List(ctx context.Context, limit, offset int) ([]*entities.User, error) {
	var users []*entities.User
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&users).Error
	return users, err
}

// generateID generates a unique ID for users
func generateID() string {
	randBytes := make([]byte, 16)
	if _, err := rand.Read(randBytes); err != nil {
		// Handle error, e.g., return a default or panic
		return "user_default_id"
	}
	return "user_" + hex.EncodeToString(randBytes)
}
