# Clean Architecture Flow - Step by Step

## 🚀 Application Startup

```
main.go
    ↓
Initialize Logger (pkg/logger)
    ↓
Load Configuration (configs/config.go)
    ↓
Create App Context (internal/app/app.go)
    ↓
Initialize Database (internal/infrastructure/database/database.go)
    ↓
Run Migrations (GORM AutoMigrate)
    ↓
Initialize Dependencies:
  ├── UserRepository (PostgreSQL implementation)
  ├── UserUseCase (Business logic)
  ├── UserHandler (HTTP handlers)
  └── Router (Chi router setup)
    ↓
Start HTTP Server
    ↓
Wait for Requests...
```

## 📡 HTTP Request Flow

### Example: Create User Request

```
1. HTTP Request: POST /api/v1/users
   {
     "email": "user@example.com",
     "name": "John Doe"
   }
    ↓
2. Chi Router (internal/interfaces/http/router/router.go)
    ↓
3. Middleware Stack:
   ├── RequestID
   ├── RealIP
   ├── Logger
   ├── Recoverer
   ├── LoggingMiddleware
   └── CORS
    ↓
4. UserHandler.CreateUser (internal/interfaces/http/handlers/user_handlers.go)
    ↓
5. UserUseCase.CreateUser (internal/usecase/user_usecase.go)
    ↓
6. Business Logic:
   ├── Validate input (email, name required)
   ├── Check if user exists
   ├── Create User entity
   └── Call repository
    ↓
7. UserRepository.Create (internal/infrastructure/database/postgres_user_repository.go)
    ↓
8. GORM Operation (pkg/postgres)
    ↓
9. PostgreSQL Database
    ↓
10. Return User Entity
    ↓
11. JSON Response:
    {
      "status": "success",
      "message": "User created successfully",
      "data": {
        "id": "user_1234567890",
        "email": "user@example.com",
        "name": "John Doe",
        "created_at": "2023-01-01T00:00:00Z",
        "updated_at": "2023-01-01T00:00:00Z"
      },
      "timestamp": "2023-01-01T00:00:00Z"
    }
```

## 🏗️ Clean Architecture Layers

```
┌─────────────────────────────────────────────────────────────┐
│                    EXTERNAL LAYER                          │
│  (Frameworks & Drivers)                                   │
├─────────────────────────────────────────────────────────────┤
│ HTTP Server (Chi) │ PostgreSQL │ GORM │ Environment Vars │
└─────────────────────────────────────────────────────────────┘
                              ↑
                              │
┌─────────────────────────────────────────────────────────────┐
│                INTERFACE ADAPTERS LAYER                   │
│  (Interface Adapters)                                     │
├─────────────────────────────────────────────────────────────┤
│ HTTP Handlers │ Router │ Middleware │ Postgres Repository │
└─────────────────────────────────────────────────────────────┘
                              ↑
                              │
┌─────────────────────────────────────────────────────────────┐
│              APPLICATION BUSINESS RULES                    │
│  (Use Cases)                                             │
├─────────────────────────────────────────────────────────────┤
│ UserUseCase │ Business Logic │ Validation │ Orchestration │
└─────────────────────────────────────────────────────────────┘
                              ↑
                              │
┌─────────────────────────────────────────────────────────────┐
│              ENTERPRISE BUSINESS RULES                    │
│  (Entities & Repository Interfaces)                      │
├─────────────────────────────────────────────────────────────┤
│ User Entity │ Repository Interface │ Domain Models │ Rules │
└─────────────────────────────────────────────────────────────┘
```

## 🔄 Dependency Flow

```
Domain Layer (Entities & Interfaces)
    ↑
Use Case Layer (Business Logic)
    ↑
Interface Layer (HTTP Handlers)
    ↑
Infrastructure Layer (Database Implementation)
```

## 📁 File Structure with Flow

```
cmd/server/main.go
    ↓ (calls)
internal/app/app.go
    ↓ (initializes)
├── configs/config.go (Load configuration)
├── internal/infrastructure/database/database.go (DB connection)
├── internal/domain/entities/user.go (User entity)
├── internal/domain/repositories/user_repository.go (Interface)
├── internal/infrastructure/database/postgres_user_repository.go (Implementation)
├── internal/usecase/user_usecase.go (Business logic)
├── internal/interfaces/http/handlers/user_handlers.go (HTTP handling)
└── internal/interfaces/http/router/router.go (Routing)
```

## 🎯 Key Principles in Action

### 1. **Dependency Inversion**
```go
// Domain Layer defines the interface
type UserRepository interface {
    Create(ctx context.Context, user *entities.User) error
    GetByID(ctx context.Context, id string) (*entities.User, error)
}

// Infrastructure Layer implements the interface
type PostgresUserRepository struct {
    db *gorm.DB
}

// Use Case depends on interface, not implementation
type UserUseCase struct {
    userRepo repositories.UserRepository  // Interface, not concrete type
}
```

### 2. **Separation of Concerns**
```go
// Handler only handles HTTP
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    // Parse HTTP request
    // Call use case
    // Return HTTP response
}

// Use Case only handles business logic
func (uc *UserUseCase) CreateUser(ctx context.Context, email, name string) (*entities.User, error) {
    // Business validation
    // Business rules
    // Call repository
}

// Repository only handles data access
func (r *PostgresUserRepository) Create(ctx context.Context, user *entities.User) error {
    // Database operations only
}
```

### 3. **Single Responsibility**
- **Entities**: Pure business objects (no dependencies)
- **Use Cases**: Business logic and orchestration
- **Handlers**: HTTP request/response handling
- **Repositories**: Data access operations
- **Middleware**: Cross-cutting concerns

## 🔧 Configuration Flow

```
Environment Variables
    ↓
envconfig Library
    ↓
Config Struct
    ↓
Database Connection
    ↓
Application Startup
```

### Environment Variables → Config
```bash
# Environment Variables
DATABASE_HOST=localhost
DATABASE_USER=postgres
DATABASE_PASSWORD=password
DATABASE_DBNAME=jackpot
SERVER_PORT=8080
LOG_LEVEL=info

# → Config Struct
type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
    Log      LogConfig
}
```

## 🧪 Testing Strategy

```
Unit Tests
├── Use Case Tests (with mock repository)
├── Handler Tests (with mock use case)
└── Repository Tests (with test database)

Integration Tests
├── End-to-End API Tests
└── Database Integration Tests

Acceptance Tests
└── Full Stack Tests
```

## 🚀 Deployment Flow

```
Source Code
    ↓
Go Build
    ↓
Binary
    ↓
Docker Build
    ↓
Container Image
    ↓
Environment Variables
    ↓
Database Connection
    ↓
Application Startup
    ↓
Health Check
    ↓
Ready for Traffic
```

## 💡 Benefits Demonstrated

1. **Testability**: Each layer can be tested independently
2. **Maintainability**: Changes in one layer don't affect others
3. **Flexibility**: Easy to swap implementations (e.g., different databases)
4. **Scalability**: Easy to add new features
5. **Independence**: Business logic is framework-agnostic

## 🔍 Real Example: Create User Flow

```go
// 1. HTTP Request comes in
POST /api/v1/users
{
  "email": "user@example.com",
  "name": "John Doe"
}

// 2. Router routes to handler
r.Post("/", userHandler.CreateUser)

// 3. Handler processes HTTP
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    var req CreateUserRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        // Return HTTP error
    }
    
    user, err := h.userUseCase.CreateUser(r.Context(), req.Email, req.Name)
    if err != nil {
        // Return HTTP error
    }
    
    // Return success response
}

// 4. Use Case handles business logic
func (uc *UserUseCase) CreateUser(ctx context.Context, email, name string) (*entities.User, error) {
    // Validate input
    if email == "" {
        return nil, errors.New("email is required")
    }
    
    // Check if user exists
    existingUser, err := uc.userRepo.GetByEmail(ctx, email)
    if err == nil && existingUser != nil {
        return nil, errors.New("user with this email already exists")
    }
    
    // Create user entity
    user := entities.NewUser(email, name)
    
    // Save to repository
    err = uc.userRepo.Create(ctx, user)
    if err != nil {
        return nil, err
    }
    
    return user, nil
}

// 5. Repository handles data access
func (r *PostgresUserRepository) Create(ctx context.Context, user *entities.User) error {
    return r.db.WithContext(ctx).Create(user).Error
}

// 6. GORM handles SQL
INSERT INTO users (id, email, name, created_at, updated_at) 
VALUES ('user_1234567890', 'user@example.com', 'John Doe', '2023-01-01T00:00:00Z', '2023-01-01T00:00:00Z');

// 7. Response flows back up
{
  "status": "success",
  "message": "User created successfully",
  "data": {
    "id": "user_1234567890",
    "email": "user@example.com",
    "name": "John Doe",
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-01T00:00:00Z"
  },
  "timestamp": "2023-01-01T00:00:00Z"
}
```

This flow demonstrates how Clean Architecture maintains separation of concerns while providing a clear, testable, and maintainable codebase. 