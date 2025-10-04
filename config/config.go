package config

import "os"

// DSN returns the PostgreSQL connection string.
// Example: postgres://user:password@localhost:5432/reminder_db?sslmode=disable
func DSN() string {
	if v := os.Getenv("DATABASE_URL"); v != "" {
		return v
	}
	// fallback for local dev
	return "host=localhost user=postgres password=admin dbname=reminder_db port=5432 sslmode=disable"
}

// HTTPPort returns the server port (default 8080).
func HTTPPort() string {
	if v := os.Getenv("PORT"); v != "" {
		return v
	}
	return "8082"
}
