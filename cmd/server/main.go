package main

import (
	"log"

	"styleai-backend/internal/config"
	"styleai-backend/internal/database"
	"styleai-backend/internal/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	cfg := config.LoadConfig()

	database.EnsureDatabaseExists(cfg)
	database.ConnectDB(cfg)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{cfg.Server.FrontendURL},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	routes.RegisterRoutes(r, cfg)

	log.Println("Server running on port", cfg.Server.Port)

	r.Run(":8080")
}
