package user

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id UserID) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByPhoneNumber(ctx context.Context, phoneNumber string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id UserID) error
	WithTx(tx *gorm.DB) Repository
}
