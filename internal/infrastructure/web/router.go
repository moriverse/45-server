package web

import (
	"github.com/gin-gonic/gin"

	"github.com/moriverse/45-server/internal/infrastructure/config"
	"github.com/moriverse/45-server/internal/infrastructure/web/handler"
	"github.com/moriverse/45-server/internal/infrastructure/web/middleware"
)

func NewRouter(
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	mw *middleware.Middleware,
	cfg config.Config,
) *gin.Engine {
	router := gin.Default()

	// Middlewares
	router.Use(mw.LoggingMiddleware())

	// Public routes
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// API v1 routes
	apiV1 := router.Group("/api/v1")
	{
		// Public auth routes
		authRoutes := apiV1.Group("/auth")
		authHandler.RegisterRoutes(authRoutes)

		// Private routes requiring authentication
		privateRoutes := apiV1.Group("/")
		privateRoutes.Use(mw.AuthMiddleware())
		{
			userHandler.RegisterRoutes(privateRoutes)
		}
	}

	return router
}
