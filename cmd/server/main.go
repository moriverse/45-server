package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/moriverse/45-server/internal/app/auth"
	"github.com/moriverse/45-server/internal/app/user"
	"github.com/moriverse/45-server/internal/infrastructure/cache"
	"github.com/moriverse/45-server/internal/infrastructure/config"
	"github.com/moriverse/45-server/internal/infrastructure/logger"
	"github.com/moriverse/45-server/internal/infrastructure/persistence"
	"github.com/moriverse/45-server/internal/infrastructure/persistence/repository"
	"github.com/moriverse/45-server/internal/infrastructure/web"
	"github.com/moriverse/45-server/internal/infrastructure/web/handler"
	"github.com/moriverse/45-server/internal/infrastructure/web/middleware"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("./configs")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	appLogger := logger.NewLogger(cfg.Log)
	appLogger.Info("Logger initialized")

	// Initialize the application
	app, err := InitializeApp(cfg, appLogger)
	if err != nil {
		appLogger.Error("Failed to initialize application", "error", err)
		os.Exit(1)
	}

	// Start the server
	appLogger.Info("Starting server", "port", cfg.Server.Port)
	if err := app.Run(":" + cfg.Server.Port); err != nil {
		appLogger.Error("Failed to run server", "error", err)
		os.Exit(1)
	}
}

func InitializeApp(cfg config.Config, appLogger *slog.Logger) (*gin.Engine, error) {
	db, err := persistence.NewDB(cfg.Database)
	if err != nil {
		return nil, err
	}

	redisClient := cache.NewRedisClient(cfg.Redis)

	userRepo := repository.NewUserRepository(db)
	authRepo := repository.NewAuthRepository(db)

	uow := persistence.NewUnitOfWork(db, userRepo, authRepo)

	// Initialize services
	authService := auth.NewService(uow, cfg.JWT)
	activityService := user.NewActivityService(redisClient, userRepo, appLogger)
	_ = user.NewService(userRepo) // Main user service

	// Initialize handlers and middleware
	authHandler := handler.NewAuthHandler(authService)
	mw := middleware.NewMiddleware(activityService, cfg.JWT, appLogger)

	return web.NewRouter(authHandler, mw, cfg), nil
}
