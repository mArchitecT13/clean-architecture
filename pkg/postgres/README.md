# PostgreSQL GORM Package

A lightweight, reusable Go package for PostgreSQL database connections using GORM. This package is designed to be portable and can be easily integrated into any Go project.

## Features

- Simple and clean API
- GORM-based PostgreSQL connection
- Connection pooling, keepalive, max open/idle connections, and connection lifetime
- Build DSN from connection fields (host, user, password, dbname, etc.)
- No dependencies on specific project structures
- Easy to integrate into any Go project
- Proper connection management with close functionality

## Installation

Add the required dependencies to your `go.mod`:

```go
require (
    gorm.io/gorm v1.30.1
    gorm.io/driver/postgres v1.6.0
)
```

## Usage

### Build DSN from Connection Fields

```go
import "your-project/pkg/postgres"

opts := postgres.ConnectionOptions{
    Host:     "localhost",
    Port:     5432,
    User:     "myuser",
    Password: "mypassword",
    DBName:   "mydb",
    SSLMode:  "disable",
    Params: map[string]string{
        "timezone": "UTC",
    },
}
dsn := postgres.BuildDSN(opts)
```

### Basic Connection with Pooling and Lifetime Options

```go
import (
    "log"
    "time"
    "your-project/pkg/postgres"
)

func main() {
    opts := postgres.ConnectionOptions{
        Host:     "localhost",
        Port:     5432,
        User:     "myuser",
        Password: "mypassword",
        DBName:   "mydb",
        SSLMode:  "disable",
    }
    dsn := postgres.BuildDSN(opts)
    config := postgres.Config{
        DSN:             dsn,
        MaxOpenConns:     20,                // Maximum number of open connections
        MaxIdleConns:     10,                // Maximum number of idle connections
        ConnMaxLifetime:  30 * time.Minute,  // Maximum connection lifetime
        ConnMaxIdleTime:  5 * time.Minute,   // Maximum idle time
    }
    db, err := postgres.New(config)
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    defer postgres.Close(db)

    // Use the database connection
    // db is a *gorm.DB instance
}
```

### Connection String Formats

The package supports various PostgreSQL connection string formats:

```go
// Key-value format
dsn := "host=localhost user=myuser password=mypassword dbname=mydb port=5432 sslmode=disable"

// URL format
dsn := "postgres://user:password@localhost:5432/dbname?sslmode=disable"

// Environment-based format
dsn := "user=myuser password=mypassword dbname=mydb host=localhost port=5432"
```

### Error Handling

Always check for errors when creating connections:

```go
config := postgres.Config{DSN: dsn}
db, err := postgres.New(config)
if err != nil {
    // Handle connection error
    log.Fatal("Database connection failed:", err)
}
defer postgres.Close(db)
```

### Using with GORM Models

```go
type User struct {
    ID    uint   `gorm:"primaryKey"`
    Name  string `gorm:"size:255"`
    Email string `gorm:"uniqueIndex"`
}

func main() {
    opts := postgres.ConnectionOptions{
        Host:     "localhost",
        Port:     5432,
        User:     "myuser",
        Password: "mypassword",
        DBName:   "mydb",
        SSLMode:  "disable",
    }
    dsn := postgres.BuildDSN(opts)
    config := postgres.Config{DSN: dsn}
    db, err := postgres.New(config)
    if err != nil {
        log.Fatal(err)
    }
    defer postgres.Close(db)

    // Auto migrate your models
    db.AutoMigrate(&User{})

    // Create a new user
    user := User{Name: "John Doe", Email: "john@example.com"}
    db.Create(&user)

    // Query users
    var users []User
    db.Find(&users)
}
```

## API Reference

### Types

#### `type ConnectionOptions struct`

| Field    | Type              | Description                  |
|----------|-------------------|------------------------------|
| Host     | string            | Database host                |
| Port     | int               | Database port                |
| User     | string            | Username                     |
| Password | string            | Password                     |
| DBName   | string            | Database name                |
| SSLMode  | string            | SSL mode (disable, require)  |
| Params   | map[string]string | Additional query parameters  |

#### `type Config struct`

| Field                | Type          | Description                                      |
|----------------------|---------------|--------------------------------------------------|
| DSN                  | string        | PostgreSQL connection string                      |
| MaxOpenConns         | int           | Maximum number of open connections                |
| MaxIdleConns         | int           | Maximum number of idle connections                |
| ConnMaxLifetime      | time.Duration | Maximum time a connection may be reused           |
| ConnMaxIdleTime      | time.Duration | Maximum time a connection may be idle             |

### Functions

#### `BuildDSN(opts ConnectionOptions) string`
Builds a DSN string from connection options.

#### `New(config Config) (*gorm.DB, error)`
Creates a new GORM database connection to PostgreSQL with advanced pooling and lifetime options.

#### `Close(db *gorm.DB) error`
Closes the database connection.

## Integration with Other Projects

This package is designed to be easily copied or imported into other Go projects. To use it in a new project:

1. Copy the `pkg/postgres` directory to your project
2. Add the required GORM dependencies to your `go.mod`
3. Import and use the package as shown in the examples above

## Dependencies

- `gorm.io/gorm` - GORM ORM library
- `gorm.io/driver/postgres` - PostgreSQL driver for GORM

## Testing

Run the tests to ensure the package works correctly:

```bash
go test ./pkg/postgres
```

## License

This package is part of the clean-architecture project and follows the same license terms.

## Contributing

When contributing to this package:

1. Ensure all tests pass
2. Follow Go best practices
3. Keep the package lightweight and focused
4. Maintain backward compatibility
5. Add appropriate documentation for new features 