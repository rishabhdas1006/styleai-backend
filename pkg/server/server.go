package server

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"styleai-backend/internal/config"
	"styleai-backend/internal/database"
	"styleai-backend/internal/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	initOnce sync.Once
	initErr  error
	app      http.Handler
)

func Handler() (http.Handler, error) {
	initOnce.Do(func() {
		cfg := config.LoadConfig()

		if err := database.Init(cfg); err != nil {
			initErr = fmt.Errorf("database initialization failed: %w", err)
			return
		}

		switch cfg.Environment {
		case "production":
			gin.SetMode(gin.ReleaseMode)
		case "testing":
			gin.SetMode(gin.TestMode)
		default:
			gin.SetMode(gin.DebugMode)
		}

		r := gin.Default()
		r.Use(cors.New(cors.Config{
			AllowOrigins:     []string{cfg.Server.FrontendURL},
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}))

		routes.RegisterRoutes(r, cfg)
		app = r
	})

	if initErr != nil {
		return nil, initErr
	}

	return app, nil
}
