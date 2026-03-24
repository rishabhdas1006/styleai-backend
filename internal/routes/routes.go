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

	// Auth
	auth := r.Group("/auth")
	{
		auth.POST("/register", userHandler.Register)
		auth.POST("/login", userHandler.Login)
	}

	// Protected user routes
	user := r.Group("/user")
	user.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
	{
		user.GET("/profile", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "protected route"})
		})
	}

	// Public category routes
	categories := r.Group("/categories")
	{
		categories.GET("/:id", categoryHandler.GetCategoryByID)
	}

	// Public product routes
	products := r.Group("/products")
	{
		products.GET("", productHandler.GetProducts)
		products.GET("/:id", productHandler.GetProductByID)
		products.GET("/:id/variants", variantHandler.GetVariants)
	}

	// Admin routes
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
