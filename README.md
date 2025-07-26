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
│   ├── postgres/
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
The application uses environment variables for configuration. You can set them directly or create a `.env` file.

#### Available Environment Variables:

**Server Configuration:**
- `SERVER_HOST` - Server host (default: localhost)
- `SERVER_PORT` - Server port (default: 8080)

**Database Configuration:**
- `DATABASE_HOST` - Database host (default: localhost)
- `DATABASE_PORT` - Database port (default: 5432)
- `DATABASE_USER` - Database user (default: postgres)
- `DATABASE_PASSWORD` - Database password (default: password)
- `DATABASE_DBNAME` - Database name (default: jackpot)
- `DATABASE_SSLMODE` - SSL mode (default: disable)
- `DATABASE_MAX_OPEN_CONNS` - Max open connections (default: 20)
- `DATABASE_MAX_IDLE_CONNS` - Max idle connections (default: 10)
- `DATABASE_CONN_MAX_LIFETIME` - Connection max lifetime (default: 30m)
- `DATABASE_CONN_MAX_IDLE_TIME` - Connection max idle time (default: 5m)

**Logging Configuration:**
- `LOG_LEVEL` - Log level (default: info)

#### Example Usage:
```bash
# Set environment variables directly
export DATABASE_USER=postgres
export DATABASE_PASSWORD=mypassword
export DATABASE_DBNAME=myapp
go run cmd/server/main.go

# Or use a .env file (copy env.example to .env and modify)
cp env.example .env
# Edit .env with your values
go run cmd/server/main.go
```

## Key Features

- **Clean Architecture**: Strict separation of concerns
- **Chi Router**: Lightweight and fast HTTP router
- **Structured Logging**: Using logrus with context
- **App Context**: Centralized application state management
- **Middleware Support**: CORS, logging, authentication ready
- **Testable**: Easy to unit test with dependency injection
- **Reusable Packages**: Modular packages in `pkg/` for use in other projects

## Development

### Adding New Features
1. Define entities in `internal/domain/entities/`
2. Create repository interfaces in `internal/domain/repositories/`
3. Implement business logic in `internal/usecase/`
4. Add HTTP handlers in `internal/interfaces/http/handlers/`
5. Wire everything in `internal/app/app.go`

### Reusable Packages

The project includes several reusable packages in the `pkg/` directory:

#### PostgreSQL Package (`pkg/postgres/`)
A lightweight, reusable Go package for PostgreSQL database connections using GORM.

```go
import "your-project/pkg/postgres"

// Create a new database connection
dsn := "host=localhost user=myuser password=mypassword dbname=mydb port=5432 sslmode=disable"
db, err := postgres.New(dsn)
if err != nil {
    log.Fatal("Failed to connect to database:", err)
}
defer postgres.Close(db)
```

For detailed usage and API reference, see [pkg/postgres/README.md](pkg/postgres/README.md).

### Testing
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...
```

## API Documentation

After running the server, visit [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html) for interactive API docs.

To regenerate the docs after changing handler annotations, run:

```bash
make swag
```

This will update the `docs/` directory with the latest Swagger documentation.

## Project Dependencies

- **Chi**: HTTP router and middleware
- **Logrus**: Structured logging
- **Testify**: Testing utilities
- **GORM**: ORM for database operations
- **PostgreSQL Driver**: Database driver for PostgreSQL
- **Envconfig**: Environment variable configuration

## License

MIT License 