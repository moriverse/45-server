package router

import (
	"github.com/gin-gonic/gin"

	"github.com/moriverse/45-server/internal/app/auth"
	"github.com/moriverse/45-server/internal/infrastructure/config"
	"github.com/moriverse/45-server/internal/infrastructure/web/middleware"
)

func NewRouter(authService *auth.Service, cfg config.Config) *gin.Engine {
	router := gin.Default()

	// Public routes
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// Auth routes
	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/register", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "register placeholder"})
		})
		authRoutes.POST("/login", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "login placeholder"})
		})
	}

	// Private route group
	v1 := router.Group("/api/v1")
	v1.Use(middleware.AuthMiddleware(cfg.JWT))
	{
		// User routes will be added here
	}

	return router
}
