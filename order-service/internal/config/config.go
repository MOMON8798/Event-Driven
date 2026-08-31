package config

import "os"

type Config struct {
	DBDSN    string
	HTTPPort string
}

func Load() *Config {
	return &Config{
		DBDSN:    getEnv("DB_DSN", "postgres://user:password@localhost:5432/dbname?sslmode=disable"),
		HTTPPort: getEnv("HTTP_PORT", "8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
