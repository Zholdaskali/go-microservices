package repository

import (
	"auth-service/internal/domain"
	"context"
	"errors"
)

var (
	ErrUserExists = errors.New("User Exists exception")
	ErrNotFound   = errors.New("User Not Found exception")
)

type UserRepository interface {
	// CRUD методы
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id string) error

	// TODO дальше query реализовать
}

type RefreshTokenRepository interface {
	Create(ctx context.Context, refreshtoken *domain.RefreshToken)
	GetByToken(ctx context.Context, token string) (*domain.RefreshToken, error)
	DeleteByToken(ctx context.Context, token string) error
}
