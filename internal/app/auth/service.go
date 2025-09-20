package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/moriverse/45-server/internal/domain/auth"
	domainErrors "github.com/moriverse/45-server/internal/domain/errors"
	"github.com/moriverse/45-server/internal/domain/unitofwork"
	"github.com/moriverse/45-server/internal/domain/user"
	"github.com/moriverse/45-server/internal/infrastructure/config"
	"github.com/moriverse/45-server/internal/infrastructure/wechat"
	"github.com/moriverse/45-server/internal/utils"
)

// AuthService defines the interface for authentication-related operations.
type AuthService interface {
	LoginOrRegisterWithWechat(
		ctx context.Context,
		params LoginOrRegisterWithWechatParams,
	) (*LoginOrRegisterResult, error)
}

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
) AuthService {
	return &Service{
		uow:          uow,
		jwtConfig:    jwtConfig,
		wechatClient: wechatClient,
	}
}

// LoginOrRegisterResult contains the result of a successful user login or registration.
type LoginOrRegisterResult struct {
	User    *user.User
	Token   string
	NewUser bool
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
) (*LoginOrRegisterResult, error) {
	openID, err := s.wechatClient.CodeToOpenID(ctx, params.Code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange wechat code: %w", err)
	}

	var u *user.User
	isNewUser := false
	err = s.uow.Execute(ctx, func(repos unitofwork.Repositories) error {
		existingAuth, err := repos.Auths.FindByProvider(ctx, auth.Wechat, openID)
		if err != nil {
			return err
		}

		if existingAuth != nil {
			foundUser, err := repos.Users.FindByID(ctx, existingAuth.UserID)
			if err != nil {
				return err
			}
			if foundUser == nil {
				return domainErrors.DataInconsistencyError{
					Details: fmt.Sprintf("auth record '%s' found but user '%s' is missing", existingAuth.ID, existingAuth.UserID),
				}
			}
			u = foundUser
			return nil
		}

		isNewUser = true
		now := time.Now()
		createdUser := &user.User{
			ID:        user.UserID(uuid.New().String()),
			Source:    params.Source,
			CreatedAt: now,
			UpdatedAt: now,
		}
		if err := repos.Users.Create(ctx, createdUser); err != nil {
			return err
		}

		newAuth := &auth.Auth{
			ID:         auth.AuthID(uuid.New().String()),
			UserID:     createdUser.ID,
			Provider:   auth.Wechat,
			ProviderID: openID,
			CreatedAt:  now,
			UpdatedAt:  now,
		}
		if err := repos.Auths.Create(ctx, newAuth); err != nil {
			return err
		}

		u = createdUser
		return nil
	})

	if err != nil {
		return nil, err
	}

	token, err := utils.GenerateToken(
		string(u.ID),
		s.jwtConfig.SecretKey,
		s.jwtConfig.ExpiresInHours,
	)
	if err != nil {
		return nil, err
	}

	return &LoginOrRegisterResult{User: u, Token: token, NewUser: isNewUser}, nil
}
