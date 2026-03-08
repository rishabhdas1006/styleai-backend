package main

import (
	"log"

	"styleai-backend/internal/config"
	"styleai-backend/internal/database"
	"styleai-backend/internal/routes"

	"github.com/gin-gonic/gin"
)

func main() {

	cfg := config.LoadConfig()

	database.ConnectDB(cfg)

	r := gin.Default()

	routes.RegisterRoutes(r, cfg)

	log.Println("Server running on port", cfg.Server.Port)

	r.Run(":8080")
}
