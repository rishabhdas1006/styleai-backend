package routes

import (
	"styleai-backend/internal/config"
	"styleai-backend/internal/handler"
	"styleai-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterCartRoutes(r *gin.RouterGroup, handler *handler.CartHandler, cfg *config.Config) {
	cart := r.Group("/cart")
	cart.Use(middleware.AuthMiddleware(cfg.JWT.Secret))

	cart.POST("/items", handler.AddItem)
	cart.GET("", handler.GetCart)
	cart.PUT("/items/:id", handler.UpdateItem)
	cart.DELETE("/items/:id", handler.DeleteItem)
}
