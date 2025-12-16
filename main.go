package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"github.com/ruangsawala/backend/config"
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

	// Set Gin mode based on environment
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize Gin router
	router := gin.Default()

	// Basic health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "RuangSawala backend is running",
			"env":     cfg.Env,
		})
	})

	// Start server
	log.Printf("Starting server on port %s...", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
