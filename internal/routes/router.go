package routes

import (
	"styleai-backend/internal/config"
	"styleai-backend/internal/database"
	"styleai-backend/internal/handler"
	"styleai-backend/internal/middleware"
	"styleai-backend/internal/repository"
	"styleai-backend/internal/service"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, cfg *config.Config) {

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "StyleAI backend running"})
	})

	// Dependencies
	userRepo := repository.NewUserRepository(database.DB)
	userService := service.NewUserService(userRepo, cfg.JWT.Secret)
	userHandler := handler.NewUserHandler(userService)

	categoryRepo := repository.NewCategoryRepository(database.DB)
	categoryService := service.NewCategoryService(categoryRepo)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	productRepo := repository.NewProductRepository(database.DB)
	productService := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productService)

	variantRepo := repository.NewVariantRepository(database.DB)
	variantService := service.NewVariantService(variantRepo, productRepo)
	variantHandler := handler.NewVariantHandler(variantService)

	cartRepo := repository.NewCartRepository(database.DB)
	cartService := service.NewCartService(cartRepo, variantRepo)
	cartHandler := handler.NewCartHandler(cartService)

	api := r.Group("/api/v1")

	RegisterAuthRoutes(api, userHandler, cfg)
	RegisterProductRoutes(api, productHandler, variantHandler)
	RegisterCartRoutes(api, cartHandler, cfg)
	RegisterAdminRoutes(api, categoryHandler, productHandler, variantHandler, cfg)

	// Protected user routes
	user := r.Group("/user")
	user.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
	{
		user.GET("/profile", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "protected route"})
		})
	}
}
