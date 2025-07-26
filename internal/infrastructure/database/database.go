package database

import (
	"log"

	"clean-architecture/configs"
	"clean-architecture/internal/domain/entities"
	"clean-architecture/pkg/postgres"

	"gorm.io/gorm"
)

var db *gorm.DB

// InitDatabase initializes the PostgreSQL database connection
func InitDatabase(cfg *configs.Config) error {
	opts := postgres.ConnectionOptions{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
		SSLMode:  cfg.Database.SSLMode,
	}
	dsn := postgres.BuildDSN(opts)

	config := postgres.Config{
		DSN:             dsn,
		MaxOpenConns:    cfg.Database.MaxOpenConns,
		MaxIdleConns:    cfg.Database.MaxIdleConns,
		ConnMaxLifetime: cfg.Database.ConnMaxLifetime,
		ConnMaxIdleTime: cfg.Database.ConnMaxIdleTime,
	}

	var err error
	db, err = postgres.New(config)
	if err != nil {
		return err
	}

	log.Println("Database connection established successfully")

	// Run migrations
	if err := MigrateDatabase(); err != nil {
		return err
	}

	return nil
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return db
}

// CloseDatabase closes the database connection
func CloseDatabase() error {
	if db != nil {
		return postgres.Close(db)
	}
	return nil
}

// MigrateDatabase runs database migrations
func MigrateDatabase() error {
	if db == nil {
		return nil
	}

	// Run migrations for all entities
	if err := db.AutoMigrate(&entities.User{}); err != nil {
		return err
	}

	log.Println("Database migrations completed successfully")
	return nil
}
