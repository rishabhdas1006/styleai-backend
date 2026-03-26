package routes

import (
	"styleai-backend/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterProductRoutes(r *gin.RouterGroup, productHandler *handler.ProductHandler, variantHandler *handler.VariantHandler) {
	products := r.Group("/products")
	{
		products.GET("", productHandler.GetProducts)
		products.GET("/:id", productHandler.GetProductByID)
		products.GET("/:id/variants", variantHandler.GetVariants)
	}
}
