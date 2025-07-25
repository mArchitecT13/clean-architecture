# Clean Architecture Go Project

A Go project implementing Clean Architecture principles using Chi router.

## Project Structure

```
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── app/
│   │   └── app.go
│   ├── domain/
│   │   ├── entities/
│   │   └── repositories/
│   ├── usecase/
│   │   └── usecase.go
│   ├── interfaces/
│   │   └── http/
│   │       ├── handlers/
│   │       ├── middleware/
│   │       └── router.go
│   └── infrastructure/
│       ├── database/
│       └── external/
├── pkg/
│   ├── logger/
│   └── utils/
├── configs/
├── docs/
└── scripts/
```

## Clean Architecture Layers

### 1. Domain Layer (`internal/domain/`)
- **Entities**: Core business objects
- **Repositories**: Interface definitions for data access

### 2. Use Case Layer (`internal/usecase/`)
- Business logic and application rules
- Orchestrates domain entities and repositories

### 3. interfaces Layer (`internal/interfaces/`)
- **HTTP Handlers**: Request/response handling
- **Middleware**: Cross-cutting concerns
- **Router**: Route definitions

### 4. Infrastructure Layer (`internal/infrastructure/`)
- Database implementations
- External service integrations
- Framework-specific code

## Getting Started

### Prerequisites
- Go 1.21 or higher

### Installation
```bash
# Clone the repository
git clone <repository-url>
cd clean-architecture

# Install dependencies
go mod tidy

# Run the application
go run cmd/server/main.go
```

### Environment Variables
Create a `.env` file in the root directory:
```env
PORT=8080
LOG_LEVEL=info
```

## Key Features

- **Clean Architecture**: Strict separation of concerns
- **Chi Router**: Lightweight and fast HTTP router
- **Structured Logging**: Using logrus with context
- **App Context**: Centralized application state management
- **Middleware Support**: CORS, logging, authentication ready
- **Testable**: Easy to unit test with dependency injection

## Development

### Adding New Features
1. Define entities in `internal/domain/entities/`
2. Create repository interfaces in `internal/domain/repositories/`
3. Implement business logic in `internal/usecase/`
4. Add HTTP handlers in `internal/interfaces/http/handlers/`
5. Wire everything in `internal/app/app.go`

### Testing
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...
```

## Project Dependencies

- **Chi**: HTTP router and middleware
- **Logrus**: Structured logging
- **Testify**: Testing utilities

## License

MIT License 