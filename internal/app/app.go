package app

import (
	"context"
	"net/http"

	"clean-architecture/internal/domain/repositories"
	"clean-architecture/internal/infrastructure/database"
	"clean-architecture/internal/interfaces/http/handlers"
	"clean-architecture/internal/interfaces/http/router"
	"clean-architecture/internal/usecase"
	"clean-architecture/pkg/logger"
)

// App represents the main application context
type App struct {
	Logger logger.Logger
	Router http.Handler
	ctx    context.Context

	// Dependencies
	UserRepository repositories.UserRepository
	UserUseCase    *usecase.UserUseCase
	UserHandler    *handlers.UserHandler
}

// NewApp creates a new application instance
func NewApp(logger logger.Logger) *App {
	ctx := context.Background()

	// Initialize repositories
	userRepo := database.NewMockUserRepository()

	// Initialize use cases
	userUseCase := usecase.NewUserUseCase(userRepo, logger)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userUseCase)

	// Create router with dependencies
	r := router.NewRouter(logger, userHandler)

	return &App{
		Logger:         logger,
		Router:         r,
		ctx:            ctx,
		UserRepository: userRepo,
		UserUseCase:    userUseCase,
		UserHandler:    userHandler,
	}
}

// Context returns the application context
func (a *App) Context() context.Context {
	return a.ctx
}

// WithContext returns a new app instance with the given context
func (a *App) WithContext(ctx context.Context) *App {
	return &App{
		Logger:         a.Logger.WithContext(ctx),
		Router:         a.Router,
		ctx:            ctx,
		UserRepository: a.UserRepository,
		UserUseCase:    a.UserUseCase,
		UserHandler:    a.UserHandler,
	}
}

// Shutdown gracefully shuts down the application
func (a *App) Shutdown(ctx context.Context) error {
	a.Logger.Info("Shutting down application...")
	// Add any cleanup logic here
	return nil
}
