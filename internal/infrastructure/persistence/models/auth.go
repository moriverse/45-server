package models

import (
	"time"
)

// Auth is the persistence model for the auths table.
type Auth struct {
	ID             string    `gorm:"primaryKey;type:uuid"`
	UserID         string    `gorm:"column:user_id;type:uuid"`
	Provider       string    `gorm:"column:provider"`
	ProviderUserID string    `gorm:"column:provider_user_id"`
	PasswordHash   string    `gorm:"column:password_hash"`
	CreatedAt      time.Time `gorm:"column:created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at"`
}

func (Auth) TableName() string {
	return "auths"
}
