package persistence

import (
	"context"

	"gorm.io/gorm"

	"github.com/moriverse/45-server/internal/domain/auth"
	"github.com/moriverse/45-server/internal/domain/unitofwork"
	"github.com/moriverse/45-server/internal/domain/user"
	"github.com/moriverse/45-server/internal/infrastructure/persistence/repository"
)

// gormUnitOfWork is the GORM implementation of the UnitOfWork interface.
type gormUnitOfWork struct {
	db       *gorm.DB
	userRepo user.Repository
	authRepo auth.Repository
}

// NewUnitOfWork creates a new GORM UnitOfWork.
func NewUnitOfWork(db *gorm.DB) unitofwork.UnitOfWork {
	return &gormUnitOfWork{
		db:       db,
		userRepo: repository.NewUserRepository(db),
		authRepo: repository.NewAuthRepository(db),
	}
}

// Execute runs the given function in a single database transaction.
func (uow *gormUnitOfWork) Execute(
	ctx context.Context,
	fn func(repos unitofwork.Repositories) error,
) error {
	return uow.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		repos := unitofwork.Repositories{
			Users: repository.NewUserRepository(tx),
			Auths: repository.NewAuthRepository(tx),
		}
		return fn(repos)
	})
}

// Users returns a non-transactional user repository.
func (uow *gormUnitOfWork) Users() user.Repository {
	return uow.userRepo
}

// Auths returns a non-transactional auth repository.
func (uow *gormUnitOfWork) Auths() auth.Repository {
	return uow.authRepo
}
