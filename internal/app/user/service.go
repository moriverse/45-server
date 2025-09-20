package user

import (
	"context"

	domainErrors "github.com/moriverse/45-server/internal/domain/errors"
	"github.com/moriverse/45-server/internal/domain/user"
	"github.com/moriverse/45-server/internal/domain/unitofwork"
)

// UserService defines the interface for user-related operations.
type UserService interface {
	GetUserProfile(ctx context.Context, userID user.UserID) (*user.User, error)
	UpdateUserProfile(
		ctx context.Context,
		userID user.UserID,
		params user.UpdateProfileParams,
	) (*user.User, error)
}

// Service is the application service for user-related operations.
type Service struct {
	uow unitofwork.UnitOfWork
}

// NewService creates a new instance of the user service.
func NewService(uow unitofwork.UnitOfWork) UserService {
	return &Service{uow: uow}
}

// GetUserProfile retrieves a user's profile by their ID using a non-transactional read.
func (s *Service) GetUserProfile(ctx context.Context, userID user.UserID) (*user.User, error) {
	u, err := s.uow.Users().FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, domainErrors.NotFoundError{Resource: "user", ID: string(userID)}
	}

	return u, nil
}

// UpdateUserProfile updates a user's profile within a transaction.
func (s *Service) UpdateUserProfile(
	ctx context.Context,
	userID user.UserID,
	params user.UpdateProfileParams,
) (*user.User, error) {
	var updatedUser *user.User
	err := s.uow.Execute(ctx, func(repos unitofwork.Repositories) error {
		u, err := repos.Users.FindByID(ctx, userID)
		if err != nil {
			return err
		}
		if u == nil {
			return domainErrors.NotFoundError{Resource: "user", ID: string(userID)}
		}

		u.Update(params)

		if err := repos.Users.Update(ctx, u); err != nil {
			return err
		}
		updatedUser = u
		return nil
	})

	return updatedUser, err
}
