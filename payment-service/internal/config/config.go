package config

import "os"

type Config struct {
	DBDSN    string
	HTTPPort string
}

func Load() *Config {
	return &Config{
		DBDSN:    getEnv("PAYMENT_DB_DSN", "postgres://user:password@localhost:5433/dbname?sslmode=disable"),
		HTTPPort: getEnv("PAYMENT_HTTP_PORT", "8081"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
