package app

import (
	"context"
	"net/http"

	"clean-architecture/configs"
	"clean-architecture/internal/domain/entities"
	"clean-architecture/internal/domain/repositories"
	"clean-architecture/internal/infrastructure/database"
	"clean-architecture/internal/interfaces/http/handlers"
	"clean-architecture/internal/interfaces/http/router"
	"clean-architecture/internal/usecase"
	"clean-architecture/pkg/logger"

	"gorm.io/gorm"
)

// App represents the main application context
type App struct {
	Logger logger.Logger
	Router http.Handler
	ctx    context.Context
	DB     *gorm.DB
	Config *configs.Config

	// Dependencies
	UserRepository repositories.UserRepository
	UserUseCase    *usecase.UserUseCase
	UserHandler    *handlers.UserHandler
}

// NewApp creates a new application instance
func NewApp(logger logger.Logger) *App {
	ctx := context.Background()

	// Load configuration
	cfg, err := configs.Load()
	if err != nil {
		logger.Fatal("Failed to load configuration:", err)
	}

	// Initialize database
	if err := database.InitDatabase(cfg); err != nil {
		logger.Fatal("Failed to initialize database:", err)
	}

	// Get database instance
	db := database.GetDB()

	// Run migrations
	if err := db.AutoMigrate(&entities.User{}); err != nil {
		logger.Fatal("Failed to run database migrations:", err)
	}
	logger.Info("Database migrations completed successfully")

	// Initialize repositories
	userRepo := database.NewPostgresUserRepository(db)

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
		DB:             db,
		Config:         cfg,
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
		DB:             a.DB,
		Config:         a.Config,
		UserRepository: a.UserRepository,
		UserUseCase:    a.UserUseCase,
		UserHandler:    a.UserHandler,
	}
}

// Shutdown gracefully shuts down the application
func (a *App) Shutdown(ctx context.Context) error {
	a.Logger.Info("Shutting down application...")

	// Close database connection
	if err := database.CloseDatabase(); err != nil {
		a.Logger.Error("Failed to close database connection:", err)
	}

	return nil
}
