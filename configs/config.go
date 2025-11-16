package configs

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort string

	DBHost string
	DBUser string
	DBPass string
	DBName string
	DBPort string
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		AppPort: os.Getenv("APP_PORT"),
		DBHost:  os.Getenv("DB_HOST"),
		DBUser:  os.Getenv("DB_USER"),
		DBPass:  os.Getenv("DB_PASS"),
		DBName:  os.Getenv("DB_NAME"),
		DBPort:  os.Getenv("DB_PORT"),
	}

	return cfg, nil
}
