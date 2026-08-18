package routes

import (
	"styleai-backend/internal/handler"
	"styleai-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterProductRoutes(
	r *gin.RouterGroup,
	productHandler *handler.ProductHandler,
	variantHandler *handler.VariantHandler,
	jwtSecret string,
) {
	products := r.Group("/products")
	products.Use(middleware.OptionalAuthMiddleware(jwtSecret))

	{
		products.GET("", productHandler.GetProducts)
		products.GET("/:id", productHandler.GetProductByID)
		products.GET("/:id/variants", variantHandler.GetVariants)
	}
}
