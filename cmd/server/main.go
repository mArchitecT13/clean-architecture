package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"clean-architecture/internal/app"
	"clean-architecture/pkg/logger"
)

func main() {
	// Initialize logger
	logger := logger.New()

	// Create application context
	appCtx := app.NewApp(logger)

	// Create HTTP server using configuration
	serverAddr := fmt.Sprintf("%s:%s", appCtx.Config.Server.Host, appCtx.Config.Server.Port)
	server := &http.Server{
		Addr:    serverAddr,
		Handler: appCtx.Router,
	}

	// Start server in a goroutine
	go func() {
		logger.Info("Starting server on " + serverAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server error: " + err.Error())
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown: " + err.Error())
	}

	// Shutdown application
	if err := appCtx.Shutdown(ctx); err != nil {
		logger.Error("Application shutdown error: " + err.Error())
	}

	logger.Info("Server exited")
}
