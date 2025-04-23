package store

import (
	"fmt"

	"github.com/EngenMe/go-api-dod/config"
	"github.com/EngenMe/go-api-dod/internal/data/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// PostgresStore implements the database store using PostgreSQL
type PostgresStore struct {
	DB *gorm.DB
}

// NewPostgresStore creates a new PostgreSQL store
func NewPostgresStore(cfg config.DatabaseConfig) (*PostgresStore, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	db, err := gorm.Open(
		postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &PostgresStore{
		DB: db,
	}, nil
}

// RunMigrations runs database migrations using GORM's AutoMigrate
func (s *PostgresStore) RunMigrations() error {
	return s.DB.AutoMigrate(&models.User{})
}
