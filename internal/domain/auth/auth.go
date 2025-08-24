package auth

import (
	"time"

	"github.com/moriverse/45-server/internal/domain/user"
)

type AuthID string

type Provider string

const (
	Email    Provider = "email"
	Phone    Provider = "phone"
	Wechat   Provider = "wechat"
	Google   Provider = "google"
)

type Auth struct {
	ID             AuthID
	UserID         user.UserID
	Provider       Provider
	ProviderUserID string
	PasswordHash   string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
