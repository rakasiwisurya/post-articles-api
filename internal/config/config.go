package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all environment-driven settings for the service.
type Config struct {
	AppPort    string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

// Load reads configuration from environment variables, falling back to an
// optional .env file and sensible local defaults.
func Load() Config {
	_ = godotenv.Load()

	return Config{
		AppPort:    getEnv("APP_PORT", "8080"),
		DBHost:     getEnv("DB_HOST", "127.0.0.1"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "article"),
	}
}

// DSN builds the MySQL connection string for the application database.
func (c Config) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local&multiStatements=true",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}

// ServerDSN builds a connection string without a database selected, used to
// create the database itself when it does not exist yet.
func (c Config) ServerDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/", c.DBUser, c.DBPassword, c.DBHost, c.DBPort)
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
