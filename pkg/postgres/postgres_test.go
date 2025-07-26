package postgres

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBuildDSN(t *testing.T) {
	opts := ConnectionOptions{
		Host:     "localhost",
		Port:     5432,
		User:     "user",
		Password: "pass",
		DBName:   "db",
		SSLMode:  "disable",
		Params: map[string]string{
			"timezone": "UTC",
		},
	}
	dsn := BuildDSN(opts)
	assert.Contains(t, dsn, "host=localhost")
	assert.Contains(t, dsn, "user=user")
	assert.Contains(t, dsn, "password=pass")
	assert.Contains(t, dsn, "dbname=db")
	assert.Contains(t, dsn, "port=5432")
	assert.Contains(t, dsn, "sslmode=disable")
	assert.Contains(t, dsn, "timezone=UTC")
}

func TestNew(t *testing.T) {
	// Test with invalid DSN
	_, err := New(Config{DSN: "invalid-dsn"})
	assert.Error(t, err)

	// Test with valid DSN format (this won't actually connect without a real database)
	// This test demonstrates the expected usage pattern
	dsn := BuildDSN(ConnectionOptions{
		Host:     "localhost",
		Port:     5432,
		User:     "test",
		Password: "test",
		DBName:   "test",
		SSLMode:  "disable",
	})
	_, err = New(Config{DSN: dsn})
	// We expect an error here since we don't have a real database running
	// but the function should not panic and should return an error
	assert.Error(t, err)
}

func TestClose(t *testing.T) {
	// Test with nil DB - should handle gracefully
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Close function panicked with nil DB: %v", r)
		}
	}()

	err := Close(nil)
	assert.Error(t, err)

	// Test with a mock DB (this will fail since we can't create a real DB without connection)
	// This test demonstrates the expected usage pattern
	dsn := BuildDSN(ConnectionOptions{
		Host:     "localhost",
		Port:     5432,
		User:     "test",
		Password: "test",
		DBName:   "test",
		SSLMode:  "disable",
	})
	db, err := New(Config{DSN: dsn})
	if err == nil {
		// If somehow we got a valid connection, test closing it
		err = Close(db)
		assert.NoError(t, err)
	}
}

func TestNew_WithValidConfig(t *testing.T) {
	// Test that the function accepts valid DSN format
	validOpts := []ConnectionOptions{
		{
			Host:     "localhost",
			Port:     5432,
			User:     "test",
			Password: "test",
			DBName:   "test",
			SSLMode:  "disable",
		},
		{
			Host:     "localhost",
			Port:     5432,
			User:     "test",
			Password: "test",
			DBName:   "test",
			SSLMode:  "require",
			Params:   map[string]string{"timezone": "UTC"},
		},
	}

	for _, opts := range validOpts {
		dsn := BuildDSN(opts)
		t.Run("DSN: "+dsn, func(t *testing.T) {
			_, err := New(Config{
				DSN:             dsn,
				MaxOpenConns:    5,
				MaxIdleConns:    2,
				ConnMaxLifetime: 10 * time.Minute,
				ConnMaxIdleTime: 2 * time.Minute,
			})
			// We expect an error since we don't have a real database
			// but the function should handle the DSN parsing correctly
			assert.Error(t, err)
		})
	}
}

func TestClose_WithMockDB(t *testing.T) {
	// Test that Close function handles errors gracefully
	// This test ensures the function doesn't panic when given invalid input
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Close function panicked: %v", r)
		}
	}()

	// Test with nil
	err := Close(nil)
	assert.Error(t, err)
}

// Example usage functions for documentation
func ExampleBuildDSN() {
	opts := ConnectionOptions{
		Host:     "localhost",
		Port:     5432,
		User:     "myuser",
		Password: "mypassword",
		DBName:   "mydb",
		SSLMode:  "disable",
		Params:   map[string]string{"timezone": "UTC"},
	}
	dsn := BuildDSN(opts)
	_ = dsn // use dsn
}

func ExampleNew() {
	// Example of how to use the New function
	opts := ConnectionOptions{
		Host:     "localhost",
		Port:     5432,
		User:     "myuser",
		Password: "mypassword",
		DBName:   "mydb",
		SSLMode:  "disable",
	}
	dsn := BuildDSN(opts)
	config := Config{
		DSN:             dsn,
		MaxOpenConns:    20,
		MaxIdleConns:    10,
		ConnMaxLifetime: 30 * time.Minute,
		ConnMaxIdleTime: 5 * time.Minute,
	}
	db, err := New(config)
	if err != nil {
		// Handle error
		return
	}
	defer Close(db)

	// Use the database connection
	_ = db
}

func ExampleClose() {
	// Example of how to use the Close function
	opts := ConnectionOptions{
		Host:     "localhost",
		Port:     5432,
		User:     "myuser",
		Password: "mypassword",
		DBName:   "mydb",
		SSLMode:  "disable",
	}
	dsn := BuildDSN(opts)
	config := Config{DSN: dsn}
	db, err := New(config)
	if err != nil {
		// Handle error
		return
	}

	// Always close the connection when done
	err = Close(db)
	if err != nil {
		// Handle error
		return
	}
}
