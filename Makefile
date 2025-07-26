.PHONY: build run test clean deps lint help

# Variables
BINARY_NAME=clean-architecture
BUILD_DIR=build
MAIN_PATH=cmd/server/main.go

# Default target
all: build

# Build the application
build:
	@echo "Building application..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Run the application
run:
	@echo "Running application..."
	@go run $(MAIN_PATH)

# Run with hot reload (requires air: go install github.com/cosmtrek/air@latest)
dev:
	@echo "Running with hot reload..."
	@air

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod tidy
	@go mod download

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run linting
lint:
	@echo "Running linter..."
	@golangci-lint run

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@go clean

# Generate mocks (requires mockgen: go install github.com/golang/mock/mockgen@latest)
mocks:
	@echo "Generating mocks..."
	@mockgen -source=internal/domain/repositories/user_repository.go -destination=internal/infrastructure/database/mocks/user_repository_mock.go

# Docker build
docker-build:
	@echo "Building Docker image..."
	@docker build -t $(BINARY_NAME) .

# Docker run
docker-run:
	@echo "Running Docker container..."
	@docker run -p 8080:8080 $(BINARY_NAME)

# Help
help:
	@echo "Available targets:"
	@echo "  build         - Build the application"
	@echo "  run           - Run the application"
	@echo "  dev           - Run with hot reload (requires air)"
	@echo "  deps          - Install dependencies"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  lint          - Run linter"
	@echo "  clean         - Clean build artifacts"
	@echo "  mocks         - Generate mocks"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run Docker container"
	@echo "  help          - Show this help" 

.PHONY: swag
swag:
	swag init -g cmd/server/main.go 