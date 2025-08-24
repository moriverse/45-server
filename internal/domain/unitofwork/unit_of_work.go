package unitofwork

import (
	"context"

	"github.com/moriverse/45-server/internal/domain/auth"
	"github.com/moriverse/45-server/internal/domain/user"
)

// UserAuthWork defines the repositories that can be used in a user-auth transaction.
type UserAuthWork interface {
	Users() user.Repository
	Auths() auth.Repository
}

// UnitOfWork is an interface for managing transactional units of work.
type UnitOfWork interface {
	Execute(ctx context.Context, fn func(UserAuthWork) error) error
}
