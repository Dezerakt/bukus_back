// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"bukus/internal/entity"
	"context"
)

//go:generate mockgen -source=contracts.go -destination=./mocks_usecase_test.go -package=usecase_test

type (
	// Book -.
	Book interface {
		Store(ctx context.Context, b entity.Book) error
		List(ctx context.Context) ([]entity.Book, error)
		GetByID(ctx context.Context, id string) (entity.Book, error)
		Update(ctx context.Context, b entity.Book) error
		Delete(ctx context.Context, id string) error
	}

	// BookRepo -.
	BookRepo interface {
		Store(ctx context.Context, book entity.Book) error
		List(ctx context.Context) ([]entity.Book, error)
		GetByID(ctx context.Context, id string) (entity.Book, error)
		Update(ctx context.Context, book entity.Book) error
		Delete(ctx context.Context, id string) error
	}

	// Auth -.
	Auth interface {
		Register(ctx context.Context, email, password string) error
		Login(ctx context.Context, email, password string) (AccessToken string, RefreshToken string, err error)
		Refresh(ctx context.Context, refreshToken string) (AccessToken string, RefreshToken string, err error)
	}

	// UserRepo -.
	UserRepo interface {
		Store(ctx context.Context, user entity.User) error
		FindByEmail(ctx context.Context, email string) (entity.User, error)
	}
)
