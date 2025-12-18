package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Config struct to hold configuration variables
type Config struct {
	ServerPort  string
	DatabaseURL string
	DBDriver    string
}

// LoadConfig loads configuration from environment variables
func LoadConfig(path string) (config Config, err error) {
	err = godotenv.Load(path + "/.env")
	if err != nil {
		// If .env file doesn't exist, try to use environment variables directly
		err = nil
	}

	config = Config{
		ServerPort:  getEnv("SERVER_PORT", "3000"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		DBDriver:    getEnv("DB_DRIVER", "postgres"),
	}

	return config, nil
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
