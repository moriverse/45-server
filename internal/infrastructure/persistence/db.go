package persistence

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/moriverse/45-server/internal/infrastructure/config"
)

func NewDB(cfg config.DatabaseConfig) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{})
}
