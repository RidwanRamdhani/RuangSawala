package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ruangsawala/backend/controllers"
)

func SetupAuthRoutes(auth *gin.RouterGroup, authController *controllers.AuthController) {
	auth.POST("/register", authController.Register)
	auth.POST("/login", authController.Login)
}
