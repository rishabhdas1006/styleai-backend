package routes

import (
	"styleai-backend/internal/database"
	"styleai-backend/internal/handler"
	"styleai-backend/internal/repository"
	"styleai-backend/internal/service"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "StyleAI backend running",
		})
	})

	// Dependency chain
	userRepo := repository.NewUserRepository(database.DB)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	auth := r.Group("/auth")
	{
		auth.POST("/register", userHandler.Register)
		auth.POST("/login", userHandler.Login)
	}
}
