package database

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"clean-architecture/internal/domain/entities"
)

func TestMockUserRepository_Create(t *testing.T) {
	repo := NewMockUserRepository()

	tests := []struct {
		name    string
		user    *entities.User
		wantErr bool
	}{
		{
			name: "successful user creation",
			user: &entities.User{
				Email: "test@example.com",
				Name:  "Test User",
			},
			wantErr: false,
		},
		{
			name: "duplicate email",
			user: &entities.User{
				Email: "test@example.com",
				Name:  "Test User",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Create(context.Background(), tt.user)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, tt.user.ID)
				// Timestamps are set by the entity, so we just check they're not zero
				assert.False(t, tt.user.CreatedAt.IsZero())
				assert.False(t, tt.user.UpdatedAt.IsZero())
			}
		})
	}
}

func TestMockUserRepository_GetByID(t *testing.T) {
	repo := NewMockUserRepository()

	// Create a test user first
	user := &entities.User{
		Email: "test@example.com",
		Name:  "Test User",
	}
	err := repo.Create(context.Background(), user)
	assert.NoError(t, err)

	tests := []struct {
		name    string
		id      string
		want    *entities.User
		wantErr bool
	}{
		{
			name:    "existing user",
			id:      user.ID,
			want:    user,
			wantErr: false,
		},
		{
			name:    "non-existing user",
			id:      "non-existing-id",
			want:    nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.GetByID(context.Background(), tt.id)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.want != nil {
					assert.Equal(t, tt.want.ID, got.ID)
					assert.Equal(t, tt.want.Email, got.Email)
					assert.Equal(t, tt.want.Name, got.Name)
				} else {
					assert.Nil(t, got)
				}
			}
		})
	}
}

func TestMockUserRepository_GetByEmail(t *testing.T) {
	repo := NewMockUserRepository()

	// Create a test user first
	user := &entities.User{
		Email: "test@example.com",
		Name:  "Test User",
	}
	err := repo.Create(context.Background(), user)
	assert.NoError(t, err)

	tests := []struct {
		name    string
		email   string
		want    *entities.User
		wantErr bool
	}{
		{
			name:    "existing user",
			email:   "test@example.com",
			want:    user,
			wantErr: false,
		},
		{
			name:    "non-existing user",
			email:   "nonexistent@example.com",
			want:    nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.GetByEmail(context.Background(), tt.email)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.want != nil {
					assert.Equal(t, tt.want.ID, got.ID)
					assert.Equal(t, tt.want.Email, got.Email)
					assert.Equal(t, tt.want.Name, got.Name)
				} else {
					assert.Nil(t, got)
				}
			}
		})
	}
}

func TestMockUserRepository_Update(t *testing.T) {
	repo := NewMockUserRepository()

	// Create a test user first
	user := &entities.User{
		Email: "test@example.com",
		Name:  "Test User",
	}
	err := repo.Create(context.Background(), user)
	assert.NoError(t, err)

	// Store original timestamps
	originalCreatedAt := user.CreatedAt
	originalUpdatedAt := user.UpdatedAt

	tests := []struct {
		name    string
		user    *entities.User
		wantErr bool
	}{
		{
			name: "successful update",
			user: &entities.User{
				ID:    user.ID,
				Email: "updated@example.com",
				Name:  "Updated User",
			},
			wantErr: false,
		},
		{
			name: "non-existing user",
			user: &entities.User{
				ID:    "non-existing-id",
				Email: "updated@example.com",
				Name:  "Updated User",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Update(context.Background(), tt.user)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				// Verify the user was updated
				updatedUser, err := repo.GetByID(context.Background(), tt.user.ID)
				assert.NoError(t, err)
				assert.Equal(t, tt.user.Email, updatedUser.Email)
				assert.Equal(t, tt.user.Name, updatedUser.Name)
				assert.Equal(t, originalCreatedAt, updatedUser.CreatedAt)                                                        // CreatedAt should not change
				assert.True(t, updatedUser.UpdatedAt.After(originalUpdatedAt) || updatedUser.UpdatedAt.Equal(originalUpdatedAt)) // UpdatedAt should change or be equal
			}
		})
	}
}

func TestMockUserRepository_Delete(t *testing.T) {
	repo := NewMockUserRepository()

	// Create a test user first
	user := &entities.User{
		Email: "test@example.com",
		Name:  "Test User",
	}
	err := repo.Create(context.Background(), user)
	assert.NoError(t, err)

	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{
			name:    "successful deletion",
			id:      user.ID,
			wantErr: false,
		},
		{
			name:    "non-existing user",
			id:      "non-existing-id",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Delete(context.Background(), tt.id)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				// Verify the user was deleted
				deletedUser, err := repo.GetByID(context.Background(), tt.id)
				assert.NoError(t, err)
				assert.Nil(t, deletedUser)
			}
		})
	}
}

func TestMockUserRepository_List(t *testing.T) {
	repo := NewMockUserRepository()

	// Create multiple test users
	users := []*entities.User{
		{Email: "user1@example.com", Name: "User 1"},
		{Email: "user2@example.com", Name: "User 2"},
		{Email: "user3@example.com", Name: "User 3"},
		{Email: "user4@example.com", Name: "User 4"},
		{Email: "user5@example.com", Name: "User 5"},
	}

	for _, user := range users {
		err := repo.Create(context.Background(), user)
		assert.NoError(t, err)
	}

	tests := []struct {
		name     string
		limit    int
		offset   int
		expected int
	}{
		{
			name:     "list all users",
			limit:    10,
			offset:   0,
			expected: 5,
		},
		{
			name:     "list with limit",
			limit:    3,
			offset:   0,
			expected: 3,
		},
		{
			name:     "list with offset",
			limit:    10,
			offset:   2,
			expected: 3,
		},
		{
			name:     "list with limit and offset",
			limit:    2,
			offset:   1,
			expected: 2,
		},
		{
			name:     "empty result",
			limit:    10,
			offset:   10,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.List(context.Background(), tt.limit, tt.offset)

			assert.NoError(t, err)
			assert.Len(t, got, tt.expected)
		})
	}
}

func TestMockUserRepository_Concurrency(t *testing.T) {
	repo := NewMockUserRepository()

	// Test concurrent access
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func(id int) {
			user := &entities.User{
				Email: fmt.Sprintf("user%d@example.com", id),
				Name:  fmt.Sprintf("User %d", id),
			}

			// Create user
			err := repo.Create(context.Background(), user)
			assert.NoError(t, err)

			// Get user
			retrieved, err := repo.GetByID(context.Background(), user.ID)
			assert.NoError(t, err)
			assert.Equal(t, user.Email, retrieved.Email)

			// Update user
			user.Name = fmt.Sprintf("Updated User %d", id)
			err = repo.Update(context.Background(), user)
			assert.NoError(t, err)

			// Delete user
			err = repo.Delete(context.Background(), user.ID)
			assert.NoError(t, err)

			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestMockUserRepository_DataIntegrity(t *testing.T) {
	repo := NewMockUserRepository()

	// Test that returned users are copies, not references
	user := &entities.User{
		Email: "test@example.com",
		Name:  "Test User",
	}

	err := repo.Create(context.Background(), user)
	assert.NoError(t, err)

	// Get the user
	retrieved, err := repo.GetByID(context.Background(), user.ID)
	assert.NoError(t, err)

	// Modify the retrieved user
	originalName := retrieved.Name
	retrieved.Name = "Modified Name"

	// Get the user again to verify it wasn't modified in the repository
	retrievedAgain, err := repo.GetByID(context.Background(), user.ID)
	assert.NoError(t, err)
	assert.Equal(t, originalName, retrievedAgain.Name)
	assert.NotEqual(t, retrieved.Name, retrievedAgain.Name)
}
