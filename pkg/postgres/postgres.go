package postgres

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ConnectionOptions holds individual connection parameters for PostgreSQL.
type ConnectionOptions struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
	Params   map[string]string // Additional query parameters
}

// BuildDSN builds a DSN string from ConnectionOptions.
func BuildDSN(opts ConnectionOptions) string {
	params := url.Values{}
	for k, v := range opts.Params {
		params.Set(k, v)
	}
	base := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		opts.Host, opts.User, opts.Password, opts.DBName, opts.Port, opts.SSLMode)
	if len(params) > 0 {
		return base + " " + strings.ReplaceAll(params.Encode(), "&", " ")
	}
	return base
}

// Config holds configuration for the PostgreSQL connection and pool.
type Config struct {
	DSN             string        // Data Source Name
	MaxOpenConns    int           // Maximum number of open connections
	MaxIdleConns    int           // Maximum number of idle connections
	ConnMaxLifetime time.Duration // Maximum amount of time a connection may be reused
	ConnMaxIdleTime time.Duration // Maximum amount of time a connection may be idle
}

// New creates a new GORM DB connection to PostgreSQL using the provided config.
func New(config Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(config.DSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqldb, err := db.DB()
	if err != nil {
		return nil, err
	}

	if config.MaxOpenConns > 0 {
		sqldb.SetMaxOpenConns(config.MaxOpenConns)
	}
	if config.MaxIdleConns > 0 {
		sqldb.SetMaxIdleConns(config.MaxIdleConns)
	}
	if config.ConnMaxLifetime > 0 {
		sqldb.SetConnMaxLifetime(config.ConnMaxLifetime)
	}
	if config.ConnMaxIdleTime > 0 {
		sqldb.SetConnMaxIdleTime(config.ConnMaxIdleTime)
	}

	return db, nil
}

// Close closes the database connection. Use this on the returned *sql.DB from db.DB().
func Close(db *gorm.DB) error {
	if db == nil {
		return fmt.Errorf("cannot close nil database connection")
	}
	sqldb, err := db.DB()
	if err != nil {
		return err
	}
	return sqldb.Close()
}
