package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/moriverse/45-server/internal/app/auth"
	"github.com/moriverse/45-server/internal/app/user"
	"github.com/moriverse/45-server/internal/infrastructure/config"
	"github.com/moriverse/45-server/internal/infrastructure/persistence"
	"github.com/moriverse/45-server/internal/infrastructure/persistence/repository"
	"github.com/moriverse/45-server/internal/infrastructure/web/router"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("./configs")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize the application
	app, err := InitializeApp(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// Start the server
	if err := app.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

func InitializeApp(cfg config.Config) (*gin.Engine, error) {
	db, err := persistence.NewDB(cfg.Database)
	if err != nil {
		return nil, err
	}

	userRepo := repository.NewUserRepository(db)
	authRepo := repository.NewAuthRepository(db)

	uow := persistence.NewUnitOfWork(db, userRepo, authRepo)

	// To keep the dependency graph simple for now, we are not initializing the user service
	// as it is not being used.
	_ = user.NewService(userRepo)
	authService := auth.NewService(uow, cfg.JWT)

	return router.NewRouter(authService, cfg), nil
}