package config

import (
	"log"
	"os"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	Port       string
	Host       string
	Environment string
	JWTSecret  string
	CORSOrigin string
}

func LoadConfig() *Config {
	return &Config{
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnv("DB_PORT", "5432"),
		DBUser:      getEnv("DB_USER", "postgres"),
		DBPassword:  getEnv("DB_PASSWORD", "postgres"),
		DBName:      getEnv("DB_NAME", "tivix_performance_tracker"),
		DBSSLMode:   getEnv("DB_SSLMODE", "disable"),
		Port:        getEnv("PORT", "8080"),
		Host:        getEnv("HOST", "localhost"),
		Environment: getEnv("ENVIRONMENT", "development"),
		JWTSecret:   getRequiredEnv("JWT_SECRET"),
		CORSOrigin:  getEnv("CORS_ORIGIN", "http://localhost:5173"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getRequiredEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Required environment variable %s is not set", key)
	}
	if len(value) < 32 {
		log.Fatalf("Environment variable %s must be at least 32 characters long for security", key)
	}
	return value
}
