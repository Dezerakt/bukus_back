package user

import (
	"bukus/internal/entity"
	"bukus/pkg/postgres"
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

// Repo -.
type Repo struct {
	*postgres.Postgres
}

// New -.
func New(pg *postgres.Postgres) *Repo {
	return &Repo{pg}
}

// Store -.
func (obj *Repo) Store(ctx context.Context, user entity.User) error {
	now := time.Now()
	if user.ID == "" {
		user.ID = uuid.New().String()
	}

	sql, args, err := obj.Builder.
		Insert("users").
		Columns("id", "email", "password_hash", "created_at", "updated_at").
		Values(user.ID, user.Email, user.PasswordHash, now, now).
		ToSql()
	if err != nil {
		return fmt.Errorf("user repo - Store - Builder: %w", err)
	}

	_, err = obj.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("user repo - Store - Exec: %w", err)
	}

	return nil
}

// FindByEmail -.
func (obj *Repo) FindByEmail(ctx context.Context, email string) (entity.User, error) {
	sql, args, err := obj.Builder.
		Select("id", "email", "password_hash", "created_at", "updated_at").
		From("users").
		Where(squirrel.Eq{"email": email}).
		ToSql()
	if err != nil {
		return entity.User{}, fmt.Errorf("user repo - FindByEmail - Builder: %w", err)
	}

	var u entity.User
	err = obj.Pool.QueryRow(ctx, sql, args...).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return entity.User{}, fmt.Errorf("user repo - FindByEmail - QueryRow: %w", err)
	}

	return u, nil
}
