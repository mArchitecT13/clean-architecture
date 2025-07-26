package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"clean-architecture/internal/interfaces/http/handlers"
	"clean-architecture/internal/interfaces/http/middleware/logging"
	"clean-architecture/pkg/logger"

	httpSwagger "github.com/swaggo/http-swagger"
)

// NewRouter creates a new Chi router with middleware
func NewRouter(logger logger.Logger, userHandler *handlers.UserHandler) http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(logging.LoggerMiddleware(logger))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Health check endpoint
	r.Get("/health", handlers.HealthCheck)

	// Serve Swagger UI
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		// Root endpoint
		r.Get("/", handlers.RootHandler)

		// User routes
		r.Route("/users", func(r chi.Router) {
			r.Get("/", userHandler.ListUsers)
			r.Post("/", userHandler.CreateUser)
			r.Get("/{id}", userHandler.GetUser)
			r.Put("/{id}", userHandler.UpdateUser)
			r.Delete("/{id}", userHandler.DeleteUser)
		})
	})

	return r
}
