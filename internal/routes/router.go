package routes

import (
	"styleai-backend/internal/config"
	"styleai-backend/internal/database"
	"styleai-backend/internal/handler"
	"styleai-backend/internal/middleware"
	"styleai-backend/internal/repository"
	"styleai-backend/internal/service"
	"styleai-backend/pkg/utils"

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

	variantImageRepo := repository.NewVariantImageRepository(database.DB)
	cloudinaryManager, _ := utils.NewCloudinaryManager()

	variantRepo := repository.NewVariantRepository(database.DB)
	variantService := service.NewVariantService(variantRepo, variantImageRepo, productRepo, database.DB, cloudinaryManager)
	variantHandler := handler.NewVariantHandler(variantService)

	cartRepo := repository.NewCartRepository(database.DB)
	cartService := service.NewCartService(cartRepo, variantRepo)
	cartHandler := handler.NewCartHandler(cartService)

	orderRepo := repository.NewOrderRepository(database.DB)
	orderService := service.NewOrderService(cartRepo, orderRepo, variantRepo, database.DB)
	orderHandler := handler.NewOrderHandler(orderService)

	cloudinaryService := service.NewCloudinaryService(cloudinaryManager)
	cloudinaryHandler := handler.NewCloudinaryHandler(cloudinaryService)

	api := r.Group("/api/v1")

	RegisterAuthRoutes(api, userHandler, cfg)
	RegisterProductRoutes(api, productHandler, variantHandler, cfg.JWT.Secret)
	RegisterCategoryRoutes(api, categoryHandler)
	RegisterCartRoutes(api, cartHandler, cfg)
	RegisterAdminRoutes(api, categoryHandler, productHandler, variantHandler, cloudinaryHandler, cfg)
	RegisterOrderRoutes(api, orderHandler, cfg)

	// Protected user routes
	user := r.Group("/user")
	user.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
	{
		user.GET("/profile", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "protected route"})
		})
	}
}
