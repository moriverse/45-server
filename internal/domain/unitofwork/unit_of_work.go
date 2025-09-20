package unitofwork

import (
	"context"

	"github.com/moriverse/45-server/internal/domain/auth"
	"github.com/moriverse/45-server/internal/domain/user"
)

// Repositories is a container for all repository interfaces.
// It is used within the Execute method to provide transaction-scoped repositories.
type Repositories struct {
	Users user.Repository
	Auths auth.Repository
}

// UnitOfWork defines the interface for our data access layer. It provides
// methods for both transactional and non-transactional operations.
type UnitOfWork interface {
	// Execute manages a transactional unit of work. It begins a transaction,
	// provides a transaction-scoped Repositories container to the callback function,
	// and commits or rolls back the transaction based on the returned error.
	Execute(ctx context.Context, fn func(Repositories) error) error

	// Users returns a non-transactional user repository.
	Users() user.Repository
	// Auths returns a non-transactional auth repository.
	Auths() auth.Repository
}
