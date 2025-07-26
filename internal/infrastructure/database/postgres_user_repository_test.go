package database

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"clean-architecture/configs"
	"clean-architecture/internal/domain/entities"
)

func TestPostgresUserRepository_Create(t *testing.T) {
	// Skip if no database connection
	if testing.Short() {
		t.Skip("Skipping database tests in short mode")
	}

	// Load configuration
	cfg, err := configs.Load()
	require.NoError(t, err)

	// Initialize database for testing
	err = InitDatabase(cfg)
	require.NoError(t, err)
	defer CloseDatabase()

	db := GetDB()
	repo := NewPostgresUserRepository(db)

	// Clean up after test
	defer func() {
		db.Exec("DELETE FROM users")
	}()

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
				Name:  "Test User 2",
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
				assert.False(t, tt.user.CreatedAt.IsZero())
				assert.False(t, tt.user.UpdatedAt.IsZero())
			}
		})
	}
}

func TestPostgresUserRepository_GetByID(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database tests in short mode")
	}

	// Load configuration
	cfg, err := configs.Load()
	require.NoError(t, err)

	err = InitDatabase(cfg)
	require.NoError(t, err)
	defer CloseDatabase()

	db := GetDB()
	repo := NewPostgresUserRepository(db)

	defer func() {
		db.Exec("DELETE FROM users")
	}()

	// Create a test user first
	user := &entities.User{
		Email: "test@example.com",
		Name:  "Test User",
	}
	err = repo.Create(context.Background(), user)
	require.NoError(t, err)

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

func TestPostgresUserRepository_GetByEmail(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database tests in short mode")
	}

	// Load configuration
	cfg, err := configs.Load()
	require.NoError(t, err)

	err = InitDatabase(cfg)
	require.NoError(t, err)
	defer CloseDatabase()

	db := GetDB()
	repo := NewPostgresUserRepository(db)

	defer func() {
		db.Exec("DELETE FROM users")
	}()

	// Create a test user first
	user := &entities.User{
		Email: "test@example.com",
		Name:  "Test User",
	}
	err = repo.Create(context.Background(), user)
	require.NoError(t, err)

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

func TestPostgresUserRepository_Update(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database tests in short mode")
	}

	// Load configuration
	cfg, err := configs.Load()
	require.NoError(t, err)

	err = InitDatabase(cfg)
	require.NoError(t, err)
	defer CloseDatabase()

	db := GetDB()
	repo := NewPostgresUserRepository(db)

	defer func() {
		db.Exec("DELETE FROM users")
	}()

	// Create a test user first
	user := &entities.User{
		Email: "test@example.com",
		Name:  "Test User",
	}
	err = repo.Create(context.Background(), user)
	require.NoError(t, err)

	// Wait a bit to ensure timestamps are different
	time.Sleep(10 * time.Millisecond)

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
			}
		})
	}
}

func TestPostgresUserRepository_Delete(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database tests in short mode")
	}

	// Load configuration
	cfg, err := configs.Load()
	require.NoError(t, err)

	err = InitDatabase(cfg)
	require.NoError(t, err)
	defer CloseDatabase()

	db := GetDB()
	repo := NewPostgresUserRepository(db)

	defer func() {
		db.Exec("DELETE FROM users")
	}()

	// Create a test user first
	user := &entities.User{
		Email: "test@example.com",
		Name:  "Test User",
	}
	err = repo.Create(context.Background(), user)
	require.NoError(t, err)

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

func TestPostgresUserRepository_List(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database tests in short mode")
	}

	// Load configuration
	cfg, err := configs.Load()
	require.NoError(t, err)

	err = InitDatabase(cfg)
	require.NoError(t, err)
	defer CloseDatabase()

	db := GetDB()
	repo := NewPostgresUserRepository(db)

	defer func() {
		db.Exec("DELETE FROM users")
	}()

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
		require.NoError(t, err)
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
