package configs

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	// Save original environment variables
	originalEnv := map[string]string{
		"PORT":        os.Getenv("PORT"),
		"HOST":        os.Getenv("HOST"),
		"DB_HOST":     os.Getenv("DB_HOST"),
		"DB_PORT":     os.Getenv("DB_PORT"),
		"DB_USER":     os.Getenv("DB_USER"),
		"DB_PASSWORD": os.Getenv("DB_PASSWORD"),
		"DB_NAME":     os.Getenv("DB_NAME"),
		"DB_SSLMODE":  os.Getenv("DB_SSLMODE"),
		"LOG_LEVEL":   os.Getenv("LOG_LEVEL"),
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
					Host:     "localhost",
					Port:     "5432",
					User:     "postgres",
					Password: "",
					Name:     "clean_architecture",
					SSLMode:  "disable",
				},
				Log: LogConfig{
					Level: "info",
				},
			},
		},
		{
			name: "custom values",
			envVars: map[string]string{
				"PORT":        "9090",
				"HOST":        "0.0.0.0",
				"DB_HOST":     "db.example.com",
				"DB_PORT":     "5433",
				"DB_USER":     "myuser",
				"DB_PASSWORD": "mypassword",
				"DB_NAME":     "myapp",
				"DB_SSLMODE":  "require",
				"LOG_LEVEL":   "debug",
			},
			expectedConfig: &Config{
				Server: ServerConfig{
					Port: "9090",
					Host: "0.0.0.0",
				},
				Database: DatabaseConfig{
					Host:     "db.example.com",
					Port:     "5433",
					User:     "myuser",
					Password: "mypassword",
					Name:     "myapp",
					SSLMode:  "require",
				},
				Log: LogConfig{
					Level: "debug",
				},
			},
		},
		{
			name: "partial custom values",
			envVars: map[string]string{
				"PORT":      "3000",
				"LOG_LEVEL": "warn",
			},
			expectedConfig: &Config{
				Server: ServerConfig{
					Port: "3000",
					Host: "localhost",
				},
				Database: DatabaseConfig{
					Host:     "localhost",
					Port:     "5432",
					User:     "postgres",
					Password: "",
					Name:     "clean_architecture",
					SSLMode:  "disable",
				},
				Log: LogConfig{
					Level: "warn",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			// Load configuration
			config := LoadConfig()

			// Assert
			assert.Equal(t, tt.expectedConfig.Server.Port, config.Server.Port)
			assert.Equal(t, tt.expectedConfig.Server.Host, config.Server.Host)
			assert.Equal(t, tt.expectedConfig.Database.Host, config.Database.Host)
			assert.Equal(t, tt.expectedConfig.Database.Port, config.Database.Port)
			assert.Equal(t, tt.expectedConfig.Database.User, config.Database.User)
			assert.Equal(t, tt.expectedConfig.Database.Password, config.Database.Password)
			assert.Equal(t, tt.expectedConfig.Database.Name, config.Database.Name)
			assert.Equal(t, tt.expectedConfig.Database.SSLMode, config.Database.SSLMode)
			assert.Equal(t, tt.expectedConfig.Log.Level, config.Log.Level)

			// Clear environment variables for next test
			for key := range tt.envVars {
				os.Unsetenv(key)
			}
		})
	}
}

func TestGetEnv(t *testing.T) {
	// Save original environment variable
	originalValue := os.Getenv("TEST_VAR")
	defer func() {
		if originalValue != "" {
			os.Setenv("TEST_VAR", originalValue)
		} else {
			os.Unsetenv("TEST_VAR")
		}
	}()

	tests := []struct {
		name         string
		key          string
		defaultValue string
		envValue     string
		expected     string
	}{
		{
			name:         "environment variable set",
			key:          "TEST_VAR",
			defaultValue: "default",
			envValue:     "custom",
			expected:     "custom",
		},
		{
			name:         "environment variable not set",
			key:          "TEST_VAR",
			defaultValue: "default",
			envValue:     "",
			expected:     "default",
		},
		{
			name:         "empty environment variable",
			key:          "TEST_VAR",
			defaultValue: "default",
			envValue:     "",
			expected:     "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set or unset environment variable
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
			} else {
				os.Unsetenv(tt.key)
			}

			// Test getEnv function
			result := getEnv(tt.key, tt.defaultValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetEnvInt(t *testing.T) {
	// Save original environment variable
	originalValue := os.Getenv("TEST_INT_VAR")
	defer func() {
		if originalValue != "" {
			os.Setenv("TEST_INT_VAR", originalValue)
		} else {
			os.Unsetenv("TEST_INT_VAR")
		}
	}()

	tests := []struct {
		name         string
		key          string
		defaultValue int
		envValue     string
		expected     int
	}{
		{
			name:         "valid integer",
			key:          "TEST_INT_VAR",
			defaultValue: 100,
			envValue:     "42",
			expected:     42,
		},
		{
			name:         "invalid integer",
			key:          "TEST_INT_VAR",
			defaultValue: 100,
			envValue:     "not_a_number",
			expected:     100,
		},
		{
			name:         "environment variable not set",
			key:          "TEST_INT_VAR",
			defaultValue: 100,
			envValue:     "",
			expected:     100,
		},
		{
			name:         "zero value",
			key:          "TEST_INT_VAR",
			defaultValue: 100,
			envValue:     "0",
			expected:     0,
		},
		{
			name:         "negative value",
			key:          "TEST_INT_VAR",
			defaultValue: 100,
			envValue:     "-10",
			expected:     -10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set or unset environment variable
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
			} else {
				os.Unsetenv(tt.key)
			}

			// Test getEnvInt function
			result := getEnvInt(tt.key, tt.defaultValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestConfig_StructCompliance(t *testing.T) {
	// Test that config structs are properly defined
	config := &Config{
		Server: ServerConfig{
			Port: "8080",
			Host: "localhost",
		},
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     "5432",
			User:     "postgres",
			Password: "password",
			Name:     "testdb",
			SSLMode:  "disable",
		},
		Log: LogConfig{
			Level: "info",
		},
	}

	// Test ServerConfig
	assert.Equal(t, "8080", config.Server.Port)
	assert.Equal(t, "localhost", config.Server.Host)

	// Test DatabaseConfig
	assert.Equal(t, "localhost", config.Database.Host)
	assert.Equal(t, "5432", config.Database.Port)
	assert.Equal(t, "postgres", config.Database.User)
	assert.Equal(t, "password", config.Database.Password)
	assert.Equal(t, "testdb", config.Database.Name)
	assert.Equal(t, "disable", config.Database.SSLMode)

	// Test LogConfig
	assert.Equal(t, "info", config.Log.Level)
}

func TestLoadConfig_EdgeCases(t *testing.T) {
	// Save original environment variables
	originalEnv := map[string]string{
		"PORT":        os.Getenv("PORT"),
		"HOST":        os.Getenv("HOST"),
		"DB_HOST":     os.Getenv("DB_HOST"),
		"DB_PORT":     os.Getenv("DB_PORT"),
		"DB_USER":     os.Getenv("DB_USER"),
		"DB_PASSWORD": os.Getenv("DB_PASSWORD"),
		"DB_NAME":     os.Getenv("DB_NAME"),
		"DB_SSLMODE":  os.Getenv("DB_SSLMODE"),
		"LOG_LEVEL":   os.Getenv("LOG_LEVEL"),
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

	// Test with all environment variables unset
	for key := range originalEnv {
		os.Unsetenv(key)
	}

	config := LoadConfig()
	assert.NotNil(t, config)
	assert.Equal(t, "8080", config.Server.Port)
	assert.Equal(t, "localhost", config.Server.Host)
	assert.Equal(t, "localhost", config.Database.Host)
	assert.Equal(t, "5432", config.Database.Port)
	assert.Equal(t, "postgres", config.Database.User)
	assert.Equal(t, "", config.Database.Password)
	assert.Equal(t, "clean_architecture", config.Database.Name)
	assert.Equal(t, "disable", config.Database.SSLMode)
	assert.Equal(t, "info", config.Log.Level)

	// Test with empty string values
	os.Setenv("PORT", "")
	os.Setenv("HOST", "")
	os.Setenv("DB_HOST", "")
	os.Setenv("DB_PORT", "")
	os.Setenv("DB_USER", "")
	os.Setenv("DB_PASSWORD", "")
	os.Setenv("DB_NAME", "")
	os.Setenv("DB_SSLMODE", "")
	os.Setenv("LOG_LEVEL", "")

	config = LoadConfig()
	assert.NotNil(t, config)
	assert.Equal(t, "8080", config.Server.Port)
	assert.Equal(t, "localhost", config.Server.Host)
	assert.Equal(t, "localhost", config.Database.Host)
	assert.Equal(t, "5432", config.Database.Port)
	assert.Equal(t, "postgres", config.Database.User)
	assert.Equal(t, "", config.Database.Password)
	assert.Equal(t, "clean_architecture", config.Database.Name)
	assert.Equal(t, "disable", config.Database.SSLMode)
	assert.Equal(t, "info", config.Log.Level)
}
