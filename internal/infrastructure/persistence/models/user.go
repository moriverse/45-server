package models

import (
	"time"
)

// User is the persistence model for the users table.
type User struct {
	ID           string     `gorm:"primaryKey;type:uuid"`
	Email        string     `gorm:"unique"`
	PhoneNumber  string     `gorm:"column:phone_number;unique"`
	AvatarURL    string     `gorm:"column:avatar_url"`
	Source       string     `gorm:"type:user_source"`
	OnboardedAt  *time.Time `gorm:"column:onboarded_at"`
	CreatedAt    time.Time  `gorm:"column:created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at"`
	LastActiveAt *time.Time `gorm:"column:last_active_at"`
	DeletedAt    *time.Time `gorm:"column:deleted_at"`
}

func (User) TableName() string {
	return "users"
}
