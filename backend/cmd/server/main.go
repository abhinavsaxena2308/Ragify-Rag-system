package main

import (
	"ragify-backend/internal/config"
	"ragify-backend/internal/routes"
	"ragify-backend/internal/utils"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	db, err := utils.DBConnect(cfg)
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"status": "OK",
			"app":    "RAGify Backend",
		})
	})

	// Setup routes with database connection
	routes.SetupRoutes(e, db)

	// Start server
	e.Logger.Fatal(e.Start(":" + cfg.Port))
}
