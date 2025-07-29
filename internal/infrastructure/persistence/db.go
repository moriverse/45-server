package persistence

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/moriverse/45-server/internal/infrastructure/config"
)

type DB struct {
	*gorm.DB
}

func NewDB(cfg config.DatabaseConfig) (*DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}
