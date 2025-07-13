package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/moriverse/45-server/internal/infrastructure/config"
	"github.com/moriverse/45-server/internal/infrastructure/persistence"
	"github.com/moriverse/45-server/internal/infrastructure/web/middleware"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	_, err = persistence.NewDB(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize Gin router
	router := gin.Default()

	// Public route
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// Private route group
	v1 := router.Group("/api/v1")
	v1.Use(middleware.AuthMiddleware())
	{
		// User routes will be added here
	}

	// Start the server
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
