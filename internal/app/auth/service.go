package auth

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/moriverse/45-server/internal/domain/auth"
	"github.com/moriverse/45-server/internal/domain/unitofwork"
	"github.com/moriverse/45-server/internal/domain/user"
	"github.com/moriverse/45-server/internal/infrastructure/config"
	pkgAuth "github.com/moriverse/45-server/pkg/auth"
)

// Service is the application service for authentication-related operations.
type Service struct {
	uow       unitofwork.UnitOfWork
	jwtConfig config.JWTConfig
}

// NewService creates a new instance of the auth service.
func NewService(uow unitofwork.UnitOfWork, jwtConfig config.JWTConfig) *Service {
	return &Service{uow: uow, jwtConfig: jwtConfig}
}

// RegisterWithEmailParams contains the parameters for registering a new user with email.
type RegisterWithEmailParams struct {
	Email    string
	Password string
	Source   user.Source
}

// RegisterResult contains the result of a successful user registration.
type RegisterResult struct {
	User  *user.User
	Token string
}

// RegisterWithEmail creates a new user and an auth record for them using email and password,
// and returns a JWT.
func (s *Service) RegisterWithEmail(
	ctx context.Context,
	params RegisterWithEmailParams,
) (*RegisterResult, error) {
	var newUser *user.User

	err := s.uow.Execute(ctx, func(work unitofwork.UserAuthWork) error {
		// 1. Check if user with email already exists
		existingUser, err := work.Users().FindByEmail(ctx, params.Email)
		if err != nil {
			return err
		}
		if existingUser != nil {
			return ErrUserAlreadyExists
		}

		// 2. Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword(
			[]byte(params.Password),
			bcrypt.DefaultCost,
		)
		if err != nil {
			return err
		}

		// 3. Create User and Auth domain objects
		now := time.Now()
		newUser = &user.User{
			ID:        user.UserID(uuid.New().String()),
			Email:     params.Email,
			Source:    params.Source,
			CreatedAt: now,
			UpdatedAt: now,
		}

		newAuth := &auth.Auth{
			ID:           auth.AuthID(uuid.New().String()),
			UserID:       newUser.ID,
			Provider:     auth.Email,
			ProviderID:   params.Email,
			PasswordHash: string(hashedPassword),
			CreatedAt:    now,
			UpdatedAt:    now,
		}

		// 4. Save the new user and auth records
		if err := work.Users().Create(ctx, newUser); err != nil {
			return err
		}

		return work.Auths().Create(ctx, newAuth)
	})

	if err != nil {
		return nil, err
	}

	// 5. Generate JWT
	token, err := pkgAuth.GenerateToken(
		string(newUser.ID),
		s.jwtConfig.SecretKey,
		s.jwtConfig.ExpiresInHours,
	)
	if err != nil {
		return nil, err
	}

	return &RegisterResult{User: newUser, Token: token}, nil
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

// RegisterWithWechatParams contains the parameters for registering a new user via Wechat.
type RegisterWithWechatParams struct {
	Code   string // The code from Wechat OAuth
	Source user.Source
}

// RegisterWithWechat creates a new user and an auth record for them using a Wechat OAuth code.
func (s *Service) RegisterWithWechat(
	ctx context.Context,
	params RegisterWithWechatParams,
) (*RegisterResult, error) {
	// TODO: Implement Wechat registration logic
	// 1. Exchange code for openID with Wechat API
	// 2. Check if user with openID already exists
	// 3. Create user and auth records
	// 4. Generate JWT
	return nil, errors.New("wechat registration not implemented")
}

// LoginParams contains the parameters for logging in.
type LoginParams struct {
	Email    string
	Password string
}

// LoginResult contains the result of a successful login.
type LoginResult struct {
	User  *user.User
	Token string
}

// Login authenticates a user and returns a JWT.
func (s *Service) Login(ctx context.Context, params LoginParams) (*LoginResult, error) {
	var u *user.User
	var token string

	err := s.uow.Execute(ctx, func(work unitofwork.UserAuthWork) error {
		// 1. Find the auth record by email
		a, err := work.Auths().FindByProvider(ctx, auth.Email, params.Email)
		if err != nil {
			return err
		}
		if a == nil {
			return ErrInvalidCredentials
		}

		// 2. Compare the password hash
		if err := bcrypt.CompareHashAndPassword(
			[]byte(a.PasswordHash),
			[]byte(params.Password),
		); err != nil {
			return ErrInvalidCredentials
		}

		// 3. Find the user record
		u, err = work.Users().FindByID(ctx, a.UserID)
		if err != nil {
			return err
		}
		if u == nil {
			// This should not happen if the data is consistent
			return errors.New("user not found for a valid auth record")
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 4. Generate JWT
	token, err = pkgAuth.GenerateToken(
		string(u.ID),
		s.jwtConfig.SecretKey,
		s.jwtConfig.ExpiresInHours,
	)
	if err != nil {
		return nil, err
	}

	return &LoginResult{User: u, Token: token}, nil
}
