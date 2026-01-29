package config

import (
	"log"
	"os"
	"strconv"
)

// Config holds the application configuration
type Config struct {
	Port      string
	Env       string
	DBPath    string
	RedisAddr string
	RedisPass string
	RedisDB   int
}

// Load loads configuration from environment variables
func Load() *Config {
	redisDB, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))

	cfg := &Config{
		Port:      getEnv("PORT", "8080"),
		Env:       getEnv("ENV", "development"),
		DBPath:    getEnv("DB_PATH", "db/ruangsawala.db"),
		RedisAddr: getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPass: getEnv("REDIS_PASS", ""),
		RedisDB:   redisDB,
	}

	log.Printf("Configuration loaded: Port=%s, Env=%s, DBPath=%s, RedisAddr=%s", cfg.Port, cfg.Env, cfg.DBPath, cfg.RedisAddr)
	return cfg
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
