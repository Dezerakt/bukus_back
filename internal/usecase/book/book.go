package book

import (
	"bukus/internal/entity"
	"bukus/internal/usecase"
	"context"
)

// UseCase -.
type UseCase struct {
	repo usecase.BookRepo
}

// New -.
func New(r usecase.BookRepo) *UseCase {
	return &UseCase{
		repo: r,
	}
}

// Store -.
func (obj *UseCase) Store(ctx context.Context, b entity.Book) error {
	return obj.repo.Store(ctx, b)
}

// List -.
func (obj *UseCase) List(ctx context.Context) ([]entity.Book, error) {
	return obj.repo.List(ctx)
}

// GetByID -.
func (obj *UseCase) GetByID(ctx context.Context, id string) (entity.Book, error) {
	return obj.repo.GetByID(ctx, id)
}

// Update -.
func (obj *UseCase) Update(ctx context.Context, b entity.Book) error {
	return obj.repo.Update(ctx, b)
}

// Delete -.
func (obj *UseCase) Delete(ctx context.Context, id string) error {
	return obj.repo.Delete(ctx, id)
}
