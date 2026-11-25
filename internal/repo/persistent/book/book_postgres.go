package book

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
func (obj *Repo) Store(ctx context.Context, book entity.Book) error {
	now := time.Now()
	if book.ID == "" {
		book.ID = uuid.New().String()
	}

	sql, args, err := obj.Builder.
		Insert("books").
		Columns("id", "name", "user_id", "created_at", "updated_at").
		Values(book.ID, book.Name, book.UserID, now, now).
		ToSql()
	if err != nil {
		return fmt.Errorf("book repo - Store - Builder: %w", err)
	}

	_, err = obj.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("book repo - Store - Exec: %w", err)
	}

	return nil
}

// List -.
func (obj *Repo) List(ctx context.Context) ([]entity.Book, error) {
	sql, args, err := obj.Builder.
		Select("id", "name", "user_id", "created_at", "updated_at").
		From("books").
		OrderBy("created_at DESC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("book repo - List - Builder: %w", err)
	}

	rows, err := obj.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("book repo - List - Query: %w", err)
	}
	defer rows.Close()

	var books []entity.Book
	for rows.Next() {
		var b entity.Book
		if err := rows.Scan(&b.ID, &b.Name, &b.UserID, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, fmt.Errorf("book repo - List - Scan: %w", err)
		}
		books = append(books, b)
	}

	return books, nil
}

// GetByID -.
func (obj *Repo) GetByID(ctx context.Context, id string) (entity.Book, error) {
	sql, args, err := obj.Builder.
		Select("id", "name", "user_id", "created_at", "updated_at").
		From("books").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return entity.Book{}, fmt.Errorf("book repo - GetByID - Builder: %w", err)
	}

	var b entity.Book
	err = obj.Pool.QueryRow(ctx, sql, args...).Scan(&b.ID, &b.Name, &b.UserID, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		return entity.Book{}, fmt.Errorf("book repo - GetByID - QueryRow: %w", err)
	}

	return b, nil
}

// Update -.
func (obj *Repo) Update(ctx context.Context, book entity.Book) error {
	sql, args, err := obj.Builder.
		Update("books").
		Set("name", book.Name).
		Set("updated_at", time.Now()).
		Where(squirrel.Eq{"id": book.ID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("book repo - Update - Builder: %w", err)
	}

	_, err = obj.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("book repo - Update - Exec: %w", err)
	}

	return nil
}

// Delete -.
func (obj *Repo) Delete(ctx context.Context, id string) error {
	sql, args, err := obj.Builder.
		Delete("books").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("book repo - Delete - Builder: %w", err)
	}

	_, err = obj.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("book repo - Delete - Exec: %w", err)
	}

	return nil
}
