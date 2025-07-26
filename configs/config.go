package configs

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig   `envconfig:"SERVER"`
	Database DatabaseConfig `envconfig:"DATABASE"`
	Log      LogConfig      `envconfig:"LOG"`
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port string `envconfig:"PORT" default:"8080"`
	Host string `envconfig:"HOST" default:"localhost"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host            string        `envconfig:"HOST" default:"localhost"`
	Port            int           `envconfig:"PORT" default:"5432"`
	User            string        `envconfig:"USER" default:"postgres"`
	Password        string        `envconfig:"PASSWORD" default:"password"`
	DBName          string        `envconfig:"DBNAME" default:"jackpot"`
	SSLMode         string        `envconfig:"SSLMODE" default:"disable"`
	MaxOpenConns    int           `envconfig:"MAX_OPEN_CONNS" default:"20"`
	MaxIdleConns    int           `envconfig:"MAX_IDLE_CONNS" default:"10"`
	ConnMaxLifetime time.Duration `envconfig:"CONN_MAX_LIFETIME" default:"30m"`
	ConnMaxIdleTime time.Duration `envconfig:"CONN_MAX_IDLE_TIME" default:"5m"`
}

// LogConfig holds logging configuration
type LogConfig struct {
	Level string `envconfig:"LEVEL" default:"info"`
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
