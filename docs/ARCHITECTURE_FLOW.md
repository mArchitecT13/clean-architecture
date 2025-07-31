# Clean Architecture Flow Diagram

## Overview

This document illustrates the complete flow of the Clean Architecture implementation, from application startup to HTTP request handling and database operations.

## Application Startup Flow

```mermaid
graph TD
    A[main.go] --> B[Initialize Logger]
    B --> C[Load Configuration]
    C --> D[Create App Context]
    D --> E[Initialize Database]
    E --> F[Run Migrations]
    F --> G[Initialize Repositories]
    G --> H[Initialize Use Cases]
    H --> I[Initialize Handlers]
    I --> J[Setup Router]
    J --> K[Start HTTP Server]
    K --> L[Wait for Requests]

    subgraph "Configuration Layer"
        C1[configs.Load] --> C2[envconfig.Process]
        C2 --> C3[Load Environment Variables]
        C3 --> C4[Return Config Struct]
    end

    subgraph "Database Layer"
        E1[database.InitDatabase] --> E2[pkg/postgres.New]
        E2 --> E3[GORM Connection]
        E3 --> E4[Connection Pooling]
    end

    subgraph "Dependency Injection"
        G1[NewPostgresUserRepository] --> G2[NewUserUseCase]
        G2 --> G3[NewUserHandler]
        G3 --> G4[NewRouter]
    end
```

## HTTP Request Flow

```mermaid
sequenceDiagram
    participant Client
    participant Router
    participant Middleware
    participant Handler
    participant UseCase
    participant Repository
    participant Database
    participant Logger

    Client->>Router: HTTP Request
    Router->>Middleware: Process Request
    Middleware->>Logger: Log Request
    Middleware->>Handler: Route to Handler
    Handler->>UseCase: Call Business Logic
    UseCase->>Repository: Data Operation
    Repository->>Database: SQL Query
    Database-->>Repository: Query Result
    Repository-->>UseCase: Domain Entity
    UseCase-->>Handler: Business Result
    Handler-->>Middleware: HTTP Response
    Middleware->>Logger: Log Response
    Middleware-->>Client: JSON Response
```

## Clean Architecture Layers

```mermaid
graph TB
    subgraph "External Layer (Frameworks & Drivers)"
        A1[HTTP Server]
        A2[Chi Router]
        A3[PostgreSQL]
        A4[GORM]
    end

    subgraph "Interface Adapters Layer"
        B1[HTTP Handlers]
        B2[Router Setup]
        B3[Middleware]
        B4[Postgres Repository]
    end

    subgraph "Application Business Rules"
        C1[User Use Case]
        C2[Business Logic]
        C3[Validation]
        C4[Orchestration]
    end

    subgraph "Enterprise Business Rules"
        D1[User Entity]
        D2[Repository Interface]
        D3[Domain Models]
        D4[Business Rules]
    end

    A1 --> B1
    A2 --> B2
    A3 --> B4
    A4 --> B4
    B1 --> C1
    B4 --> C1
    C1 --> D1
    C1 --> D2
    D2 --> B4
```

## Detailed Request Flow Example

### 1. Create User Request

```mermaid
graph TD
    A[POST /api/v1/users] --> B[Chi Router]
    B --> C[Middleware Stack]
    C --> D[UserHandler.CreateUser]
    D --> E[UserUseCase.CreateUser]
    E --> F[Validate Input]
    F --> G[Check Email Exists]
    G --> H[Create User Entity]
    H --> I[UserRepository.Create]
    I --> J[PostgresUserRepository]
    J --> K[GORM Create]
    K --> L[PostgreSQL]
    L --> M[Return User]
    M --> N[JSON Response]
```

### 2. Get User Request

```mermaid
graph TD
    A[GET /api/v1/users/{id}] --> B[Chi Router]
    B --> C[Middleware Stack]
    C --> D[UserHandler.GetUser]
    D --> E[UserUseCase.GetUserByID]
    E --> F[UserRepository.GetByID]
    F --> G[PostgresUserRepository]
    G --> H[GORM Query]
    H --> I[PostgreSQL]
    I --> J[Return User]
    J --> K[JSON Response]
```

## Configuration Flow

```mermaid
graph LR
    A[Environment Variables] --> B[envconfig Library]
    B --> C[Config Struct]
    C --> D[Database Connection]
    C --> E[Server Settings]
    C --> F[Logging Settings]
    
    subgraph "Environment Variables"
        A1[SERVER_HOST]
        A2[SERVER_PORT]
        A3[DATABASE_HOST]
        A4[DATABASE_USER]
        A5[DATABASE_PASSWORD]
        A6[LOG_LEVEL]
    end
    
    subgraph "Config Structure"
        C1[ServerConfig]
        C2[DatabaseConfig]
        C3[LogConfig]
    end
```

## Database Connection Flow

```mermaid
graph TD
    A[Config Load] --> B[ConnectionOptions]
    B --> C[BuildDSN]
    C --> D[PostgreSQL DSN]
    D --> E[GORM Connection]
    E --> F[Connection Pooling]
    F --> G[Database Instance]
    G --> H[Run Migrations]
    H --> I[Ready for Queries]
    
    subgraph "Connection Pool Settings"
        F1[MaxOpenConns: 20]
        F2[MaxIdleConns: 10]
        F3[ConnMaxLifetime: 30m]
        F4[ConnMaxIdleTime: 5m]
    end
```

## Error Handling Flow

```mermaid
graph TD
    A[Request] --> B{Valid Request?}
    B -->|No| C[Validation Error]
    B -->|Yes| D[Business Logic]
    D --> E{Business Rule Valid?}
    E -->|No| F[Business Error]
    E -->|Yes| G[Database Operation]
    G --> H{Database Success?}
    H -->|No| I[Database Error]
    H -->|Yes| J[Success Response]
    
    C --> K[HTTP 400]
    F --> L[HTTP 422]
    I --> M[HTTP 500]
    J --> N[HTTP 200/201]
```

## Testing Flow

```mermaid
graph TD
    A[Unit Tests] --> B[Use Case Tests]
    A --> C[Handler Tests]
    A --> D[Repository Tests]
    
    B --> E[Mock Repository]
    C --> F[Mock Use Case]
    D --> G[Test Database]
    
    H[Integration Tests] --> I[End-to-End API]
    I --> J[Real Database]
    
    K[Acceptance Tests] --> L[Full Stack]
    L --> M[Production-like Environment]
```

## Deployment Flow

```mermaid
graph TD
    A[Source Code] --> B[Go Build]
    B --> C[Binary]
    C --> D[Docker Build]
    D --> E[Container Image]
    E --> F[Environment Variables]
    F --> G[Database Connection]
    G --> H[Application Startup]
    H --> I[Health Check]
    I --> J[Ready for Traffic]
    
    subgraph "Environment Setup"
        F1[Development]
        F2[Staging]
        F3[Production]
    end
```

## Key Principles Demonstrated

### 1. Dependency Inversion
- **Domain Layer** defines interfaces
- **Infrastructure Layer** implements interfaces
- **Use Case Layer** depends on abstractions, not concretions

### 2. Separation of Concerns
- **Entities**: Pure business objects
- **Use Cases**: Application business rules
- **Interfaces**: External communication
- **Infrastructure**: Technical implementation

### 3. Single Responsibility
- Each layer has one reason to change
- Handlers only handle HTTP
- Use Cases only handle business logic
- Repositories only handle data access

### 4. Open/Closed Principle
- Easy to add new use cases without modifying existing code
- Easy to add new repositories without changing business logic
- Easy to add new handlers without changing use cases

## Benefits of This Architecture

1. **Testability**: Each layer can be tested in isolation
2. **Maintainability**: Changes in one layer don't affect others
3. **Scalability**: Easy to add new features and modify existing ones
4. **Independence**: Business logic is independent of frameworks
5. **Flexibility**: Easy to swap implementations (e.g., database, HTTP framework)

## File Structure Mapping

```
cmd/server/main.go                    # Application Entry Point
├── internal/app/app.go              # Dependency Injection Container
├── internal/domain/                 # Enterprise Business Rules
│   ├── entities/user.go            # Core Business Objects
│   └── repositories/user_repository.go # Data Access Interfaces
├── internal/usecase/               # Application Business Rules
│   └── user_usecase.go            # Business Logic
├── internal/interfaces/http/       # Interface Adapters
│   ├── handlers/user_handlers.go   # HTTP Request Handling
│   ├── middleware/logging.go       # Cross-cutting Concerns
│   └── router/router.go            # Route Definitions
├── internal/infrastructure/        # Frameworks & Drivers
│   └── database/                   # Database Implementation
│       ├── postgres_user_repository.go # PostgreSQL Implementation
│       └── database.go             # Database Connection
├── pkg/                           # Reusable Packages
│   ├── postgres/                  # PostgreSQL Package
│   ├── logger/                    # Logging Package
│   └── utils/                     # Utility Functions
└── configs/                       # Configuration Management
    └── config.go                  # Environment Configuration
```

This architecture ensures that the business logic remains independent of external concerns while providing a clear and maintainable structure for the application. 
