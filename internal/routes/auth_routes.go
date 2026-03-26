package routes

import (
	"styleai-backend/internal/config"
	"styleai-backend/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(r *gin.RouterGroup, userHandler *handler.UserHandler, cfg *config.Config) {
	auth := r.Group("/auth")
	{
		auth.POST("/register", userHandler.Register)
		auth.POST("/login", userHandler.Login)
	}
}
