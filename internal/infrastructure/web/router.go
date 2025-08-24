package web

import (
	"github.com/gin-gonic/gin"

	"github.com/moriverse/45-server/internal/infrastructure/config"
	"github.com/moriverse/45-server/internal/infrastructure/web/handler"
	"github.com/moriverse/45-server/internal/infrastructure/web/middleware"
)

func NewRouter(authHandler *handler.AuthHandler, mw *middleware.Middleware, cfg config.Config) *gin.Engine {
	router := gin.Default()

	// Middlewares
	router.Use(mw.LoggingMiddleware())

	// Public routes
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// Auth routes
	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/login", authHandler.Login)
	}

	// Private route group
	v1 := router.Group("/api/v1")
	v1.Use(mw.AuthMiddleware())
	{
		// User routes will be added here
	}

	return router
}
