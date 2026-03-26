package routes

import (
	"styleai-backend/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterCategoryRoutes(r *gin.RouterGroup, categoryHandler *handler.CategoryHandler) {
	categories := r.Group("/categories")
	{
		categories.GET("/:id", categoryHandler.GetCategoryByID)
	}
}
