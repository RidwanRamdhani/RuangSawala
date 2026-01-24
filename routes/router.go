package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ruangsawala/backend/config"
	"github.com/ruangsawala/backend/controllers"
)

func NewRouter(cfg *config.Config, authController *controllers.AuthController) *gin.Engine {
	// Set Gin mode based on environment
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Basic health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "RuangSawala backend is running",
			"env":     cfg.Env,
		})
	})

	// Auth routes
	auth := router.Group("/auth")
	SetupAuthRoutes(auth, authController)

	return router
}
