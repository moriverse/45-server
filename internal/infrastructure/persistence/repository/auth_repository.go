package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/moriverse/45-server/internal/domain/auth"
	"github.com/moriverse/45-server/internal/domain/user"
	"github.com/moriverse/45-server/internal/infrastructure/persistence/models"
)

// AuthRepository is a GORM implementation of the auth.Repository interface.
type AuthRepository struct {
	db *gorm.DB
}

// NewAuthRepository creates a new instance of AuthRepository.
func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

// WithTx returns a new instance of the repository with the database connection set to the
// given transaction.
func (r *AuthRepository) WithTx(tx *gorm.DB) auth.Repository {
	return &AuthRepository{db: tx}
}

// Create creates a new auth record in the database.
func (r *AuthRepository) Create(ctx context.Context, a *auth.Auth) error {
	model := toAuthModel(a)
	return r.db.WithContext(ctx).Create(model).Error
}

// FindByProvider finds an auth record by provider and provider user ID.
func (r *AuthRepository) FindByProvider(
	ctx context.Context,
	provider auth.Provider,
	providerID string,
) (*auth.Auth, error) {
	var model models.Auth
	if err := r.db.WithContext(ctx).First(
		&model,
		"provider = ? AND provider_id = ?",
		provider,
		providerID,
	).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Or a custom not found error
		}
		return nil, err
	}
	return toAuthDomain(&model), nil
}

// toAuthModel converts a domain auth to a GORM auth model.
func toAuthModel(a *auth.Auth) *models.Auth {
	return &models.Auth{
		ID:           string(a.ID),
		UserID:       string(a.UserID),
		Provider:     string(a.Provider),
		ProviderID:   a.ProviderID,
		PasswordHash: a.PasswordHash,
		CreatedAt:    a.CreatedAt,
		UpdatedAt:    a.UpdatedAt,
	}
}

// toAuthDomain converts a GORM auth model to a domain auth.
func toAuthDomain(m *models.Auth) *auth.Auth {
	return &auth.Auth{
		ID:           auth.AuthID(m.ID),
		UserID:       user.UserID(m.UserID),
		Provider:     auth.Provider(m.Provider),
		ProviderID:   m.ProviderID,
		PasswordHash: m.PasswordHash,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}
