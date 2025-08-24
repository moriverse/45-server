package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/moriverse/45-server/internal/domain/auth"
	"github.com/moriverse/45-server/internal/domain/unitofwork"
	"github.com/moriverse/45-server/internal/domain/user"
	"github.com/moriverse/45-server/internal/infrastructure/config"
	"github.com/moriverse/45-server/internal/infrastructure/wechat"
	"github.com/moriverse/45-server/internal/utils"
)

// Service is the application service for authentication-related operations.
type Service struct {
	uow          unitofwork.UnitOfWork
	jwtConfig    config.JWTConfig
	wechatClient *wechat.Client
}

// NewService creates a new instance of the auth service.
func NewService(
	uow unitofwork.UnitOfWork,
	jwtConfig config.JWTConfig,
	wechatClient *wechat.Client,
) *Service {
	return &Service{
		uow:          uow,
		jwtConfig:    jwtConfig,
		wechatClient: wechatClient,
	}
}

// RegisterResult contains the result of a successful user registration.
type RegisterResult struct {
	User  *user.User
	Token string
}

// RegisterWithPhoneParams contains the parameters for registering a new user with a phone number.
type RegisterWithPhoneParams struct {
	PhoneNumber string
	Code        string
	Source      user.Source
}

// RegisterWithPhone creates a new user and an auth record for them using a phone number and
// verification code.
func (s *Service) RegisterWithPhone(
	ctx context.Context,
	params RegisterWithPhoneParams,
) (*RegisterResult, error) {
	// TODO: Implement phone registration logic
	// 1. Verify the code
	// 2. Check if user with phone number already exists
	// 3. Create user and auth records
	// 4. Generate JWT
	return nil, errors.New("phone registration not implemented")
}

// LoginOrRegisterWithWechatParams contains the parameters for signing in a user via Wechat.
type LoginOrRegisterWithWechatParams struct {
	Code   string // The code from Wechat OAuth
	Source user.Source
}

// LoginOrRegisterWithWechat exchanges a Wechat code for an openid, then finds the corresponding
// user or creates a new one if they don't exist.
func (s *Service) LoginOrRegisterWithWechat(
	ctx context.Context,
	params LoginOrRegisterWithWechatParams,
) (*RegisterResult, error) {
	// 1. Exchange code for openID with Wechat API
	openID, err := s.wechatClient.CodeToOpenID(ctx, params.Code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange wechat code: %w", err)
	}

	var u *user.User
	err = s.uow.Execute(ctx, func(work unitofwork.UserAuthWork) error {
		// 2. Check if an auth record with this openID already exists
		existingAuth, err := work.Auths().FindByProvider(ctx, auth.Wechat, openID)
		if err != nil {
			return err
		}

		if existingAuth != nil {
			// User exists, so we're logging them in.
			foundUser, err := work.Users().FindByID(ctx, existingAuth.UserID)
			if err != nil {
				return err
			}
			if foundUser == nil {
				// This indicates data inconsistency and should not happen.
				return errors.New("auth record found but user is missing")
			}
			u = foundUser
			return nil
		}

		// 3. User does not exist, so we're creating them.
		now := time.Now()
		newUser := &user.User{
			ID:        user.UserID(uuid.New().String()),
			Source:    params.Source,
			CreatedAt: now,
			UpdatedAt: now,
		}
		if err := work.Users().Create(ctx, newUser); err != nil {
			return err
		}

		newAuth := &auth.Auth{
			ID:         auth.AuthID(uuid.New().String()),
			UserID:     newUser.ID,
			Provider:   auth.Wechat,
			ProviderID: openID,
			CreatedAt:  now,
			UpdatedAt:  now,
		}
		if err := work.Auths().Create(ctx, newAuth); err != nil {
			return err
		}

		u = newUser
		return nil
	})

	if err != nil {
		return nil, err
	}

	// 4. Generate JWT for the found or created user
	token, err := utils.GenerateToken(
		string(u.ID),
		s.jwtConfig.SecretKey,
		s.jwtConfig.ExpiresInHours,
	)
	if err != nil {
		return nil, err
	}

	return &RegisterResult{User: u, Token: token}, nil
}
