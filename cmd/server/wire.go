//go:build wireinject

//go:generate wire

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"

	"github.com/moriverse/45-server/internal/infrastructure/config"
	"github.com/moriverse/45-server/internal/infrastructure/persistence"
	"github.com/moriverse/45-server/internal/infrastructure/web/middleware"
)

func InitializeApp(cfg config.Config) (*gin.Engine, error) {
	wire.Build(
		persistence.NewDB,
		wire.FieldsOf(new(config.Config), "Database"),
		newRouter,
	)
	return nil, nil
}

func newRouter(db *persistence.DB, cfg config.Config) *gin.Engine {
	router := gin.Default()

	// Public route
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// Private route group
	v1 := router.Group("/api/v1")
	v1.Use(middleware.AuthMiddleware(cfg.JWT))
	{
		// User routes will be added here
	}

	return router
}
