package configs

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	// Application
	App AppConfig

	// Database
	Database DatabaseConfig

	// JWT
	JWT JWTConfig

	// CORS
	CORS CORSConfig
}

type AppConfig struct {
	Port        string
	Environment string
}

type DatabaseConfig struct {
	Host     string
	User     string
	Password string
	Name     string
	Port     string
}

type JWTConfig struct {
	Secret string
}

type CORSConfig struct {
	AllowedOrigins string
}

func Load() (*Config, error) {
	// Load .env file (optional, doesn't fail if not found)
	_ = godotenv.Load()

	// Validate required environment variables
	jwtSecret := getEnv("JWT_SECRET", "")
	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET environment variable is required and must be at least 32 characters long")
	}
	if len(jwtSecret) < 32 {
		return nil, fmt.Errorf("JWT_SECRET must be at least 32 characters long")
	}

	dbHost := getEnv("DB_HOST", "")
	dbUser := getEnv("DB_USER", "")
	dbName := getEnv("DB_NAME", "")
	dbPort := getEnv("DB_PORT", "")

	if dbHost == "" || dbUser == "" || dbName == "" || dbPort == "" {
		return nil, fmt.Errorf("missing required database environment variables (DB_HOST, DB_USER, DB_NAME, DB_PORT)")
	}

	cfg := &Config{
		App: AppConfig{
			Port:        getEnv("APP_PORT", "8080"),
			Environment: getEnv("ENV", "development"),
		},
		Database: DatabaseConfig{
			Host:     dbHost,
			User:     dbUser,
			Password: getEnv("DB_PASS", ""),
			Name:     dbName,
			Port:     dbPort,
		},
		JWT: JWTConfig{
			Secret: jwtSecret,
		},
		CORS: CORSConfig{
			AllowedOrigins: getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:5173"),
		},
	}

	return cfg, nil
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultVal string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultVal
}

// LoadConfigOrFatal loads config and logs fatal error if it fails
func LoadConfigOrFatal() *Config {
	cfg, err := Load()
	if err != nil {
		log.Fatal("❌ Failed to load config:", err)
	}
	return cfg
}
