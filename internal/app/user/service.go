package user

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/moriverse/45-server/internal/domain/user"
)

// Service is the application service for user-related operations.
type Service struct {
	userRepo user.Repository
}

// NewService creates a new instance of the user service.
func NewService(userRepo user.Repository) *Service {
	return &Service{userRepo: userRepo}
}

// CreateUserParams contains the parameters for creating a new user.
type CreateUserParams struct {
	Username    string
	Email       string
	PhoneNumber string
	AvatarURL   string
	Source      user.Source
}

// CreateUser creates a new user.
func (s *Service) CreateUser(ctx context.Context, params CreateUserParams) (*user.User, error) {
	now := time.Now()
	newUser := &user.User{
		ID:          user.UserID(uuid.New().String()),
		Username:    params.Username,
		Email:       params.Email,
		PhoneNumber: params.PhoneNumber,
		AvatarURL:   params.AvatarURL,
		Source:      params.Source,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.userRepo.Create(ctx, newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}
