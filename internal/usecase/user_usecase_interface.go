package usecase

import (
	"context"

	"clean-architecture/internal/domain/entities"
)

// UserUseCaseInterface defines the interface for user business logic
type UserUseCaseInterface interface {
	CreateUser(ctx context.Context, email, name string) (*entities.User, error)
	GetUserByID(ctx context.Context, id string) (*entities.User, error)
	UpdateUser(ctx context.Context, id, name, email string) (*entities.User, error)
	DeleteUser(ctx context.Context, id string) error
	ListUsers(ctx context.Context, limit, offset int) ([]*entities.User, error)
}
