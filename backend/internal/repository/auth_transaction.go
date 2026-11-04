package repository

import (
	"context"

	"gorm.io/gorm"
)

type AuthTransaction interface {
	Run(context.Context, func(UserRepository, RefreshTokenRepository) error) error
}

type authTransaction struct {
	db *gorm.DB
}

func NewAuthTransaction(db *gorm.DB) AuthTransaction {
	return &authTransaction{db: db}
}

func (r *authTransaction) Run(ctx context.Context, fn func(UserRepository, RefreshTokenRepository) error) error {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(NewUserRepository(tx), NewRefreshTokenRepository(tx))
	})
	return authDBError(err)
}
