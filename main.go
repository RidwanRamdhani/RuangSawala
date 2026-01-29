package main

import (
	"context"
	"database/sql"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/ruangsawala/backend/config"
	"github.com/ruangsawala/backend/controllers"
	"github.com/ruangsawala/backend/routes"
	"github.com/ruangsawala/backend/services"
	_ "modernc.org/sqlite"
)

func runMigrations(dbPath string) {
	m, err := migrate.New(
		"file://migrations",
		"sqlite3://"+dbPath,
	)
	if err != nil {
		log.Fatal("Migration init failed:", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("Migration up failed:", err)
	}
	log.Println("Migrations applied successfully")
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Load configuration
	cfg := config.Load()

	// Run database migrations
	runMigrations(cfg.DBPath)

	// Initialize database connection
	db, err := sql.Open("sqlite", cfg.DBPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Initialize Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPass,
		DB:       cfg.RedisDB,
	})
	defer rdb.Close()

	// Test Redis connection
	if err := rdb.Ping(context.TODO()).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Redis connected successfully")

	// Initialize services and controllers
	authService := services.NewAuthService(db)
	authController := controllers.NewAuthController(authService)

	// Initialize Gin router
	router := routes.NewRouter(cfg, authController)

	// Start server
	log.Printf("Starting server on port %s...", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
