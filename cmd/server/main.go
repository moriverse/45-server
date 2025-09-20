package main

import (
	"log/slog"
	"os"

	"github.com/moriverse/45-server/internal/app/auth"
	"github.com/moriverse/45-server/internal/app/user"
	"github.com/moriverse/45-server/internal/infrastructure/config"
	"github.com/moriverse/45-server/internal/infrastructure/logger"
	"github.com/moriverse/45-server/internal/infrastructure/persistence"
	"github.com/moriverse/45-server/internal/infrastructure/wechat"
	"github.com/moriverse/45-server/internal/infrastructure/web"
	"github.com/moriverse/45-server/internal/infrastructure/web/handler"
	"github.com/moriverse/45-server/internal/infrastructure/web/middleware"
)

func main() {
	// Get config path from environment variable
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		slog.Error("CONFIG_PATH environment variable is not set")
		os.Exit(1)
	}

	// Load configuration
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		slog.Error("Failed to load configuration", "path", configPath, "error", err)
		os.Exit(1)
	}

	// Initialize logger
	appLogger := logger.NewLogger(cfg.Log)
	slog.SetDefault(appLogger)

	// Initialize database connection
	db, err := persistence.NewDB(cfg.Database)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}

	// Dependency Injection
	uow := persistence.NewUnitOfWork(db)
	wechatClient := wechat.NewClient(cfg.Wechat)

	// Services
	authService := auth.NewService(uow, cfg.JWT, wechatClient)
	userService := user.NewService(uow)

	// Handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	mw := middleware.NewMiddleware(cfg.JWT, appLogger)

	// Initialize router
	router := web.NewRouter(authHandler, userHandler, mw, cfg)

	// Start server
	slog.Info("Starting server", "port", cfg.Server.Port)
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		slog.Error("Failed to start server", "error", err)
		os.Exit(1)
	}
}
