package config

import (
	"log"
	"os"
)

// Config holds the application configuration
type Config struct {
	Port string
	Env  string
}

// Load loads configuration from environment variables
func Load() *Config {
	cfg := &Config{
		Port: getEnv("PORT", "8080"),
		Env:  getEnv("ENV", "development"),
	}

	log.Printf("Configuration loaded: Port=%s, Env=%s", cfg.Port, cfg.Env)
	return cfg
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
