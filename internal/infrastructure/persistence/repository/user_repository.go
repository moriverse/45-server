package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/moriverse/45-server/internal/domain/user"
	"github.com/moriverse/45-server/internal/infrastructure/persistence/models"
)

// UserRepository is a GORM implementation of the user.Repository interface.
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new instance of UserRepository.
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// WithTx returns a new instance of the repository with the database connection set to the
// given transaction.
func (r *UserRepository) WithTx(tx *gorm.DB) user.Repository {
	return &UserRepository{db: tx}
}

// Create creates a new user in the database.
func (r *UserRepository) Create(ctx context.Context, u *user.User) error {
	model := toUserModel(u)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *UserRepository) FindByID(ctx context.Context, id user.UserID) (*user.User, error) {
	var model models.User
	if err := r.db.WithContext(ctx).First(&model, "id = ?", string(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Or a custom not found error
		}
		return nil, err
	}
	return toUserDomain(&model), nil
}

// FindByEmail finds a user by their email address.
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	var model models.User
	if err := r.db.WithContext(ctx).First(&model, "email = ?", email).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Or a custom not found error
		}
		return nil, err
	}
	return toUserDomain(&model), nil
}

// FindByPhoneNumber finds a user by their phone number.
func (r *UserRepository) FindByPhoneNumber(
	ctx context.Context,
	phoneNumber string,
) (*user.User, error) {
	var model models.User
	if err := r.db.WithContext(ctx).First(
		&model, "phone_number = ?", phoneNumber,
	).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Or a custom not found error
		}
		return nil, err
	}
	return toUserDomain(&model), nil
}

// Update updates an existing user in the database.
func (r *UserRepository) Update(ctx context.Context, u *user.User) error {
	model := toUserModel(u)
	return r.db.WithContext(ctx).Save(model).Error
}

// Delete marks a user as deleted in the database.
func (r *UserRepository) Delete(ctx context.Context, id user.UserID) error {
	return r.db.WithContext(ctx).Model(&models.User{}).
		Where("id = ?", string(id)).
		Update("deleted_at", time.Now()).Error
}

// UpdateLastActiveAt updates the last_active_at timestamp for a user.
func (r *UserRepository) UpdateLastActiveAt(ctx context.Context, id user.UserID, t time.Time) error {
	return r.db.WithContext(ctx).Model(&models.User{}).
		Where("id = ?", string(id)).
		Update("last_active_at", t).Error
}

// toUserModel converts a domain user to a GORM user model.
func toUserModel(u *user.User) *models.User {
	return &models.User{
		ID:           string(u.ID),
		Email:        u.Email,
		PhoneNumber:  u.PhoneNumber,
		AvatarURL:    u.AvatarURL,
		Source:       string(u.Source),
		OnboardedAt:  u.OnboardedAt,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
		LastActiveAt: u.LastActiveAt,
		DeletedAt:    u.DeletedAt,
	}
}

// toUserDomain converts a GORM user model to a domain user.
func toUserDomain(m *models.User) *user.User {
	return &user.User{
		ID:           user.UserID(m.ID),
		Email:        m.Email,
		PhoneNumber:  m.PhoneNumber,
		AvatarURL:    m.AvatarURL,
		Source:       user.Source(m.Source),
		OnboardedAt:  m.OnboardedAt,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
		LastActiveAt: m.LastActiveAt,
		DeletedAt:    m.DeletedAt,
	}
}
