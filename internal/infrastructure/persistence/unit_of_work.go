package persistence

import (
	"context"

	"gorm.io/gorm"

	"github.com/moriverse/45-server/internal/domain/auth"
	"github.com/moriverse/45-server/internal/domain/unitofwork"
	"github.com/moriverse/45-server/internal/domain/user"
)

// gormUnitOfWork is the GORM implementation of the UnitOfWork interface.
type gormUnitOfWork struct {
	db       *gorm.DB
	userRepo user.Repository
	authRepo auth.Repository
}

// NewUnitOfWork creates a new GORM UnitOfWork.
func NewUnitOfWork(db *gorm.DB, userRepo user.Repository, authRepo auth.Repository) unitofwork.UnitOfWork {
	return &gormUnitOfWork{
		db:       db,
		userRepo: userRepo,
		authRepo: authRepo,
	}
}

// Execute runs the given function in a single database transaction.
func (uow *gormUnitOfWork) Execute(ctx context.Context, fn func(work unitofwork.UserAuthWork) error) error {
	return uow.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		work := &gormUserAuthWork{
			userRepo: uow.userRepo.WithTx(tx),
			authRepo: uow.authRepo.WithTx(tx),
		}
		return fn(work)
	})
}

// gormUserAuthWork is the GORM implementation of the UserAuthWork interface.
type gormUserAuthWork struct {
	userRepo user.Repository
	authRepo auth.Repository
}

func (w *gormUserAuthWork) Users() user.Repository {
	return w.userRepo
}

func (w *gormUserAuthWork) Auths() auth.Repository {
	return w.authRepo
}
