package routes

import (
	"styleai-backend/internal/config"
	"styleai-backend/internal/handler"
	"styleai-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterAdminRoutes(r *gin.RouterGroup, categoryHandler *handler.CategoryHandler, productHandler *handler.ProductHandler, variantHandler *handler.VariantHandler, cfg *config.Config) {
	admin := r.Group("/admin")
	admin.Use(
		middleware.AuthMiddleware(cfg.JWT.Secret),
		middleware.AdminMiddleware(),
	)
	{
		admin.POST("/categories", categoryHandler.CreateCategory)
		admin.POST("/products", productHandler.CreateProduct)
		admin.POST("/products/:id/variants", variantHandler.CreateVariant)
		admin.PUT("/variants/:id", variantHandler.UpdateVariant)
		admin.DELETE("/variants/:id", variantHandler.DeleteVariant)
	}
}
