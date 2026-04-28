package routes

import (
	"styleai-backend/internal/config"
	"styleai-backend/internal/handler"
	"styleai-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterOrderRoutes(r *gin.RouterGroup, orderHandler *handler.OrderHandler, cfg *config.Config) {
	orders := r.Group("/orders")
	orders.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
	{
		orders.POST("/checkout", orderHandler.Checkout)
		orders.GET("", orderHandler.GetOrders)
		orders.GET("/:id", orderHandler.GetOrderByID)
	}
}
