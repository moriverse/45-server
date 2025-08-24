package auth

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, auth *Auth) error
	FindByProvider(ctx context.Context, provider Provider, providerUserID string) (*Auth, error)
	WithTx(tx *gorm.DB) Repository
}
