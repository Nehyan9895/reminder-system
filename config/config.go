package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv loads environment variables from a .env file
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, relying on system environment variables")
	}
}

// DSN returns the PostgreSQL connection string from the environment
func DSN() string {
	v := os.Getenv("DATABASE_URL")
	if v == "" {
		log.Fatal("DATABASE_URL is not set in environment or .env file")
	}
	return v
}

// HTTPPort returns the server port from environment (default 8082)
func HTTPPort() string {
	v := os.Getenv("PORT")
	if v == "" {
		return "8080"
	}
	return v
}
