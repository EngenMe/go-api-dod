package config

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Auth     AuthConfig
}

// ServerConfig holds server-specific configuration
type ServerConfig struct {
	Port int
	Mode string // development, production
}

// DatabaseConfig holds database-specific configuration
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// AuthConfig holds authentication-specific configuration
type AuthConfig struct {
	JWTSecret       string
	TokenExpiration time.Duration
	TokenIssuer     string
	BcryptCost      int
}

// Load loads the configuration from environment variables
func Load() (Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	var cfg Config

	// Server configuration
	port, err := strconv.Atoi(getEnv("SERVER_PORT", "8080"))
	if err != nil {
		return cfg, errors.New("invalid SERVER_PORT")
	}
	cfg.Server.Port = port
	cfg.Server.Mode = getEnv("SERVER_MODE", "development")

	// Database configuration
	cfg.Database.Host = getEnv("DB_HOST", "localhost")
	dbPort, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
	if err != nil {
		return cfg, errors.New("invalid DB_PORT")
	}
	cfg.Database.Port = dbPort
	cfg.Database.User = getEnv("DB_USER", "postgres")
	cfg.Database.Password = getEnv("DB_PASSWORD", "postgres")
	cfg.Database.DBName = getEnv("DB_NAME", "GO_API")
	cfg.Database.SSLMode = getEnv("DB_SSLMODE", "disable")

	// Auth configuration
	cfg.Auth.JWTSecret = getEnv("JWT_SECRET", "")
	if cfg.Auth.JWTSecret == "" {
		return cfg, errors.New("JWT_SECRET is required")
	}

	tokenExpiration, err := strconv.Atoi(getEnv("TOKEN_EXPIRATION_HOURS", "24"))
	if err != nil {
		return cfg, errors.New("invalid TOKEN_EXPIRATION_HOURS")
	}
	cfg.Auth.TokenExpiration = time.Duration(tokenExpiration) * time.Hour
	cfg.Auth.TokenIssuer = getEnv("TOKEN_ISSUER", "go-api-dod")

	bcryptCost, err := strconv.Atoi(getEnv("BCRYPT_COST", "10"))
	if err != nil {
		return cfg, errors.New("invalid BCRYPT_COST")
	}
	cfg.Auth.BcryptCost = bcryptCost

	return cfg, nil
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
