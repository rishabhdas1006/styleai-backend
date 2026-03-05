package main

import (
	"styleai/internal/database"
	"styleai/internal/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// connect database
	database.ConnectDB()

	r := gin.Default()

	// register routes
	routes.RegisterRoutes(r)

	r.Run(":8080")
}
