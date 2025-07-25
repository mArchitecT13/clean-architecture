# Clean Architecture Go Project Structure

## Overview

This is a complete Clean Architecture implementation in Go using Chi router. The project follows strict separation of concerns and dependency inversion principles.

## Project Structure

```
clean-architecture/
├── cmd/
│   └── server/
│       └── main.go                 # Application entry point
├── internal/
│   ├── app/
│   │   └── app.go                  # Main application context
│   ├── domain/
│   │   ├── entities/
│   │   │   └── user.go            # Domain entities
│   │   └── repositories/
│   │       └── user_repository.go  # Repository interfaces
│   ├── usecase/
│   │   ├── user_usecase.go        # Business logic
│   │   └── user_usecase_test.go   # Use case tests
│   ├── interfaces/
│   │   └── http/
│   │       ├── handlers/
│   │       │   ├── handlers.go     # Base handlers
│   │       │   └── user_handlers.go # User HTTP handlers
│   │       ├── middleware/
│   │       │   └── logging/
│   │       │       └── logging.go  # Logging middleware
│   │       └── router/
│   │           └── router.go       # Chi router setup
│   └── infrastructure/
│       └── database/
│           └── mock_user_repository.go # Mock repository implementation
├── pkg/
│   ├── logger/
│   │   └── logger.go               # Structured logging
│   └── utils/
│       └── response.go             # HTTP response utilities
├── configs/
│   └── config.go                   # Configuration management
├── docs/
│   └── API.md                      # API documentation
├── scripts/                        # Build and deployment scripts
├── go.mod                          # Go module file
├── go.sum                          # Go module checksums
├── README.md                       # Project documentation
├── Makefile                        # Build automation
├── Dockerfile                      # Container configuration
├── docker-compose.yml              # Development environment
├── .gitignore                      # Git ignore rules
└── STRUCTURE.md                    # This file
```

## Clean Architecture Layers

### 1. Domain Layer (`internal/domain/`)
- **Entities**: Core business objects (User)
- **Repositories**: Interface definitions for data access
- **Purpose**: Contains enterprise business rules and entities

### 2. Use Case Layer (`internal/usecase/`)
- **Business Logic**: Application-specific business rules
- **Orchestration**: Coordinates domain entities and repositories
- **Purpose**: Contains application business rules

### 3. interfaces Layer (`internal/interfaces/`)
- **HTTP Handlers**: Request/response handling
- **Middleware**: Cross-cutting concerns (logging, CORS)
- **Router**: Route definitions using Chi
- **Purpose**: Handles external interfaces (HTTP, CLI, etc.)

### 4. Infrastructure Layer (`internal/infrastructure/`)
- **Database**: Repository implementations
- **External Services**: Third-party integrations
- **Purpose**: Framework-specific implementations

## Key Features

### Application Context (`internal/app/app.go`)
- Centralized application state management
- Dependency injection container
- Logger and context propagation
- Graceful shutdown handling

### Structured Logging (`pkg/logger/`)
- JSON-formatted logs
- Context-aware logging
- Configurable log levels
- Request tracing support

### HTTP Middleware
- Request ID generation
- Real IP detection
- Structured request logging
- CORS support
- Panic recovery

### Configuration Management (`configs/`)
- Environment-based configuration
- Database settings
- Logging configuration
- Server settings

## Dependencies

### Core Dependencies
- **Chi Router**: Lightweight HTTP router
- **Logrus**: Structured logging
- **Testify**: Testing utilities

### Development Dependencies
- **Air**: Hot reload for development
- **GolangCI-Lint**: Code linting
- **Mockgen**: Mock generation

## API Endpoints

### Health & Info
- `GET /health` - Health check
- `GET /api/v1/` - API information

### User Management
- `GET /api/v1/users` - List users
- `POST /api/v1/users` - Create user
- `GET /api/v1/users/{id}` - Get user
- `PUT /api/v1/users/{id}` - Update user
- `DELETE /api/v1/users/{id}` - Delete user

## Development Workflow

### Prerequisites
- Go 1.21+
- Docker (optional)
- Make (optional)

### Quick Start
```bash
# Install dependencies
make deps

# Run the application
make run

# Run with hot reload
make dev

# Run tests
make test

# Build for production
make build
```

### Docker Development
```bash
# Start development environment
docker-compose up

# Build and run container
make docker-build
make docker-run
```

## Testing Strategy

### Unit Tests
- Use case layer testing with mock repositories
- Handler testing with mock use cases
- Repository testing with in-memory implementations

### Integration Tests
- End-to-end API testing
- Database integration testing
- Middleware testing

## Deployment

### Docker
- Multi-stage build for minimal image size
- Non-root user for security
- Health checks included
- Alpine Linux base for small footprint

### Environment Variables
- `PORT`: Server port (default: 8080)
- `LOG_LEVEL`: Logging level (default: info)
- `DB_HOST`: Database host
- `DB_PORT`: Database port
- `DB_USER`: Database user
- `DB_PASSWORD`: Database password
- `DB_NAME`: Database name

## Clean Architecture Benefits

### 1. Independence of Frameworks
- Business logic is independent of Chi router
- Easy to switch HTTP frameworks
- Database agnostic design

### 2. Testability
- Dependency injection enables easy mocking
- Business logic can be tested in isolation
- Repository pattern allows database mocking

### 3. Independence of UI
- Business logic is separate from HTTP handlers
- Easy to add CLI, gRPC, or GraphQL interfaces
- Same use cases work with different interfaces mechanisms

### 4. Independence of Database
- Repository interfaces define data access contracts
- Easy to switch between databases
- Business logic doesn't know about database details

### 5. Independence of External Agencies
- Business rules don't depend on external services
- External dependencies are injected
- Easy to mock external services for testing

## Best Practices Implemented

### 1. Dependency Inversion
- High-level modules don't depend on low-level modules
- Both depend on abstractions
- Abstractions don't depend on details

### 2. Single Responsibility
- Each layer has a single responsibility
- Use cases contain only business logic
- Handlers only handle HTTP concerns

### 3. Interface Segregation
- Repository interfaces are focused and specific
- Handlers depend only on what they need
- Clean separation between layers

### 4. Open/Closed Principle
- Easy to extend with new entities
- New use cases can be added without modifying existing code
- New interfaces mechanisms can be added

## Future Enhancements

### Planned Features
- JWT Authentication
- Rate Limiting
- Database Migrations
- API Versioning
- GraphQL Support
- gRPC Support
- Event Sourcing
- CQRS Pattern

### Monitoring & Observability
- Prometheus Metrics
- Distributed Tracing
- Structured Logging
- Health Checks
- Performance Monitoring 