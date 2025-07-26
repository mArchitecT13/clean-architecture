# 🏗️ Clean Architecture Flow - Quick Overview

## 📊 High-Level Flow

```
┌─────────────────────────────────────────────────────────────┐
│                    HTTP REQUEST                            │
│              POST /api/v1/users                           │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│                    CHI ROUTER                              │
│              internal/interfaces/http/router               │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│                  MIDDLEWARE STACK                          │
│  RequestID │ RealIP │ Logger │ Recoverer │ CORS           │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│                  HTTP HANDLER                              │
│              internal/interfaces/http/handlers             │
│              • Parse HTTP request                          │
│              • Call Use Case                               │
│              • Return HTTP response                        │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│                   USE CASE                                 │
│              internal/usecase                              │
│              • Business validation                         │
│              • Business rules                              │
│              • Call Repository                             │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│                 REPOSITORY                                 │
│              internal/infrastructure/database              │
│              • Data access operations                      │
│              • Database queries                            │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│                  DATABASE                                  │
│              PostgreSQL + GORM                             │
└─────────────────────────────────────────────────────────────┘
```

## 🎯 Clean Architecture Layers

```
┌─────────────────────────────────────────────────────────────┐
│                    EXTERNAL LAYER                          │
│  (Frameworks & Drivers)                                   │
│  HTTP Server │ Database │ GORM │ Environment Variables    │
└─────────────────────────────────────────────────────────────┘
                              ↑
┌─────────────────────────────────────────────────────────────┐
│                INTERFACE ADAPTERS                          │
│  (Interface Adapters)                                     │
│  HTTP Handlers │ Router │ Middleware │ Repository Impl    │
└─────────────────────────────────────────────────────────────┘
                              ↑
┌─────────────────────────────────────────────────────────────┐
│              APPLICATION BUSINESS RULES                    │
│  (Use Cases)                                             │
│  Business Logic │ Validation │ Orchestration              │
└─────────────────────────────────────────────────────────────┘
                              ↑
┌─────────────────────────────────────────────────────────────┐
│              ENTERPRISE BUSINESS RULES                    │
│  (Entities & Interfaces)                                 │
│  User Entity │ Repository Interface │ Domain Models       │
└─────────────────────────────────────────────────────────────┘
```

## 🔄 Dependency Direction

```
Domain Layer (Entities & Interfaces)
    ↑ (depends on)
Use Case Layer (Business Logic)
    ↑ (depends on)
Interface Layer (HTTP Handlers)
    ↑ (depends on)
Infrastructure Layer (Database Implementation)
```

## 📁 File Structure Flow

```
cmd/server/main.go
    ↓
internal/app/app.go (Dependency Injection)
    ↓
├── configs/config.go (Configuration)
├── internal/domain/entities/user.go (Business Objects)
├── internal/domain/repositories/user_repository.go (Interfaces)
├── internal/usecase/user_usecase.go (Business Logic)
├── internal/interfaces/http/handlers/user_handlers.go (HTTP)
├── internal/infrastructure/database/postgres_user_repository.go (Data)
└── internal/interfaces/http/router/router.go (Routing)
```

## 🚀 Application Startup

```
1. main.go
   ↓
2. Initialize Logger
   ↓
3. Load Configuration (Environment Variables)
   ↓
4. Create App Context
   ↓
5. Initialize Database Connection
   ↓
6. Run Database Migrations
   ↓
7. Initialize Dependencies:
   ├── UserRepository (PostgreSQL)
   ├── UserUseCase (Business Logic)
   ├── UserHandler (HTTP)
   └── Router (Chi)
   ↓
8. Start HTTP Server
   ↓
9. Wait for Requests...
```

## 📡 Request Flow Example

```
HTTP Request: POST /api/v1/users
{
  "email": "user@example.com",
  "name": "John Doe"
}
    ↓
Chi Router → Middleware → Handler
    ↓
UserHandler.CreateUser()
    ↓
UserUseCase.CreateUser()
    ↓
Business Logic:
├── Validate input
├── Check if user exists
├── Create User entity
└── Call repository
    ↓
UserRepository.Create()
    ↓
GORM → PostgreSQL
    ↓
Return User Entity
    ↓
JSON Response:
{
  "status": "success",
  "message": "User created successfully",
  "data": { ... }
}
```

## 🎯 Key Principles

### 1. **Dependency Inversion**
- Domain defines interfaces
- Infrastructure implements interfaces
- Use Cases depend on abstractions

### 2. **Separation of Concerns**
- **Entities**: Pure business objects
- **Use Cases**: Business logic
- **Handlers**: HTTP handling
- **Repositories**: Data access

### 3. **Single Responsibility**
- Each layer has one reason to change
- Each component has one job

### 4. **Open/Closed Principle**
- Easy to add new features
- Easy to modify existing features
- No breaking changes to existing code

## 💡 Benefits

✅ **Testability**: Each layer can be tested independently  
✅ **Maintainability**: Changes don't affect other layers  
✅ **Flexibility**: Easy to swap implementations  
✅ **Scalability**: Easy to add new features  
✅ **Independence**: Business logic is framework-agnostic  

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

## 🧪 Testing Strategy

```
Unit Tests
├── Use Case Tests (with mock repository)
├── Handler Tests (with mock use case)
└── Repository Tests (with test database)

Integration Tests
├── End-to-End API Tests
└── Database Integration Tests
```

This Clean Architecture implementation provides a clear, maintainable, and testable structure that follows SOLID principles and separation of concerns. 