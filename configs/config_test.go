package configs

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	// Save original environment variables
	originalEnv := map[string]string{
		"SERVER_HOST":       os.Getenv("SERVER_HOST"),
		"SERVER_PORT":       os.Getenv("SERVER_PORT"),
		"DATABASE_HOST":     os.Getenv("DATABASE_HOST"),
		"DATABASE_PORT":     os.Getenv("DATABASE_PORT"),
		"DATABASE_USER":     os.Getenv("DATABASE_USER"),
		"DATABASE_PASSWORD": os.Getenv("DATABASE_PASSWORD"),
		"DATABASE_DBNAME":   os.Getenv("DATABASE_DBNAME"),
		"DATABASE_SSLMODE":  os.Getenv("DATABASE_SSLMODE"),
		"LOG_LEVEL":         os.Getenv("LOG_LEVEL"),
	}

	// Restore environment variables after test
	defer func() {
		for key, value := range originalEnv {
			if value != "" {
				os.Setenv(key, value)
			} else {
				os.Unsetenv(key)
			}
		}
	}()

	tests := []struct {
		name           string
		envVars        map[string]string
		expectedConfig *Config
	}{
		{
			name:    "default values",
			envVars: map[string]string{},
			expectedConfig: &Config{
				Server: ServerConfig{
					Port: "8080",
					Host: "localhost",
				},
				Database: DatabaseConfig{
					Host:            "localhost",
					Port:            5432,
					User:            "postgres",
					Password:        "password",
					DBName:          "jackpot",
					SSLMode:         "disable",
					MaxOpenConns:    20,
					MaxIdleConns:    10,
					ConnMaxLifetime: 30 * time.Minute,
					ConnMaxIdleTime: 5 * time.Minute,
				},
				Log: LogConfig{
					Level: "info",
				},
			},
		},
		{
			name: "custom values",
			envVars: map[string]string{
				"SERVER_HOST":                 "0.0.0.0",
				"SERVER_PORT":                 "9090",
				"DATABASE_HOST":               "db.example.com",
				"DATABASE_PORT":               "5433",
				"DATABASE_USER":               "myuser",
				"DATABASE_PASSWORD":           "mypassword",
				"DATABASE_DBNAME":             "myapp",
				"DATABASE_SSLMODE":            "require",
				"DATABASE_MAX_OPEN_CONNS":     "50",
				"DATABASE_MAX_IDLE_CONNS":     "25",
				"DATABASE_CONN_MAX_LIFETIME":  "1h",
				"DATABASE_CONN_MAX_IDLE_TIME": "10m",
				"LOG_LEVEL":                   "debug",
			},
			expectedConfig: &Config{
				Server: ServerConfig{
					Port: "9090",
					Host: "0.0.0.0",
				},
				Database: DatabaseConfig{
					Host:            "db.example.com",
					Port:            5433,
					User:            "myuser",
					Password:        "mypassword",
					DBName:          "myapp",
					SSLMode:         "require",
					MaxOpenConns:    50,
					MaxIdleConns:    25,
					ConnMaxLifetime: time.Hour,
					ConnMaxIdleTime: 10 * time.Minute,
				},
				Log: LogConfig{
					Level: "debug",
				},
			},
		},
		{
			name: "partial custom values",
			envVars: map[string]string{
				"SERVER_PORT": "3000",
				"LOG_LEVEL":   "warn",
			},
			expectedConfig: &Config{
				Server: ServerConfig{
					Port: "3000",
					Host: "localhost",
				},
				Database: DatabaseConfig{
					Host:            "localhost",
					Port:            5432,
					User:            "postgres",
					Password:        "password",
					DBName:          "jackpot",
					SSLMode:         "disable",
					MaxOpenConns:    20,
					MaxIdleConns:    10,
					ConnMaxLifetime: 30 * time.Minute,
					ConnMaxIdleTime: 5 * time.Minute,
				},
				Log: LogConfig{
					Level: "warn",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear all environment variables first
			envVars := []string{
				"SERVER_HOST", "SERVER_PORT", "DATABASE_HOST", "DATABASE_PORT",
				"DATABASE_USER", "DATABASE_PASSWORD", "DATABASE_DBNAME", "DATABASE_SSLMODE",
				"DATABASE_MAX_OPEN_CONNS", "DATABASE_MAX_IDLE_CONNS", "DATABASE_CONN_MAX_LIFETIME",
				"DATABASE_CONN_MAX_IDLE_TIME", "LOG_LEVEL", "USER",
			}

			for _, envVar := range envVars {
				os.Unsetenv(envVar)
			}

			// Set environment variables for test
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			// Load configuration
			config, err := Load()
			assert.NoError(t, err)
			assert.NotNil(t, config)

			// Assert server config
			assert.Equal(t, tt.expectedConfig.Server.Host, config.Server.Host)
			assert.Equal(t, tt.expectedConfig.Server.Port, config.Server.Port)

			// Assert database config
			assert.Equal(t, tt.expectedConfig.Database.Host, config.Database.Host)
			assert.Equal(t, tt.expectedConfig.Database.Port, config.Database.Port)
			assert.Equal(t, tt.expectedConfig.Database.User, config.Database.User)
			assert.Equal(t, tt.expectedConfig.Database.Password, config.Database.Password)
			assert.Equal(t, tt.expectedConfig.Database.DBName, config.Database.DBName)
			assert.Equal(t, tt.expectedConfig.Database.SSLMode, config.Database.SSLMode)
			assert.Equal(t, tt.expectedConfig.Database.MaxOpenConns, config.Database.MaxOpenConns)
			assert.Equal(t, tt.expectedConfig.Database.MaxIdleConns, config.Database.MaxIdleConns)
			assert.Equal(t, tt.expectedConfig.Database.ConnMaxLifetime, config.Database.ConnMaxLifetime)
			assert.Equal(t, tt.expectedConfig.Database.ConnMaxIdleTime, config.Database.ConnMaxIdleTime)

			// Assert log config
			assert.Equal(t, tt.expectedConfig.Log.Level, config.Log.Level)
		})
	}
}

func TestLoad_EdgeCases(t *testing.T) {
	// Save original environment variables
	originalEnv := map[string]string{
		"SERVER_HOST":       os.Getenv("SERVER_HOST"),
		"SERVER_PORT":       os.Getenv("SERVER_PORT"),
		"DATABASE_HOST":     os.Getenv("DATABASE_HOST"),
		"DATABASE_PORT":     os.Getenv("DATABASE_PORT"),
		"DATABASE_USER":     os.Getenv("DATABASE_USER"),
		"DATABASE_PASSWORD": os.Getenv("DATABASE_PASSWORD"),
		"DATABASE_DBNAME":   os.Getenv("DATABASE_DBNAME"),
		"DATABASE_SSLMODE":  os.Getenv("DATABASE_SSLMODE"),
		"LOG_LEVEL":         os.Getenv("LOG_LEVEL"),
	}

	// Restore environment variables after test
	defer func() {
		for key, value := range originalEnv {
			if value != "" {
				os.Setenv(key, value)
			} else {
				os.Unsetenv(key)
			}
		}
	}()

	t.Run("invalid port number", func(t *testing.T) {
		os.Setenv("DATABASE_PORT", "invalid")

		config, err := Load()
		assert.Error(t, err)
		assert.Nil(t, config)
	})

	t.Run("invalid duration", func(t *testing.T) {
		os.Setenv("DATABASE_CONN_MAX_LIFETIME", "invalid")

		config, err := Load()
		assert.Error(t, err)
		assert.Nil(t, config)
	})

	t.Run("empty environment", func(t *testing.T) {
		// Clear all relevant environment variables
		envVars := []string{
			"SERVER_HOST", "SERVER_PORT", "DATABASE_HOST", "DATABASE_PORT",
			"DATABASE_USER", "DATABASE_PASSWORD", "DATABASE_DBNAME", "DATABASE_SSLMODE",
			"DATABASE_MAX_OPEN_CONNS", "DATABASE_MAX_IDLE_CONNS", "DATABASE_CONN_MAX_LIFETIME",
			"DATABASE_CONN_MAX_IDLE_TIME", "LOG_LEVEL", "USER",
		}

		for _, envVar := range envVars {
			os.Unsetenv(envVar)
		}

		config, err := Load()
		assert.NoError(t, err)
		assert.NotNil(t, config)

		// Should use default values
		assert.Equal(t, "localhost", config.Server.Host)
		assert.Equal(t, "8080", config.Server.Port)
		assert.Equal(t, "localhost", config.Database.Host)
		assert.Equal(t, 5432, config.Database.Port)
		assert.Equal(t, "postgres", config.Database.User)
		assert.Equal(t, "password", config.Database.Password)
		assert.Equal(t, "jackpot", config.Database.DBName)
		assert.Equal(t, "disable", config.Database.SSLMode)
		assert.Equal(t, 20, config.Database.MaxOpenConns)
		assert.Equal(t, 10, config.Database.MaxIdleConns)
		assert.Equal(t, 30*time.Minute, config.Database.ConnMaxLifetime)
		assert.Equal(t, 5*time.Minute, config.Database.ConnMaxIdleTime)
		assert.Equal(t, "info", config.Log.Level)
	})
}
