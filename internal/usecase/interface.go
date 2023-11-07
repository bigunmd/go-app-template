package usecase

import (
	"context"
	"goapptemplate/internal/domain"

	"github.com/google/uuid"
)

type (
	Books interface {
		New(ctx context.Context, book *domain.Book) (*domain.Book, error)
		View(ctx context.Context, bookID uuid.UUID) (*domain.Book, error)
		List(ctx context.Context, filters *domain.BookFilters) (*domain.BookPage, error)
		Modify(ctx context.Context, book *domain.Book) (*domain.Book, error)
		Remove(ctx context.Context, bookID uuid.UUID) error
	}
	BooksRepo interface {
		Store(ctx context.Context, book *domain.Book) (*domain.Book, error)
		Retrieve(ctx context.Context, bookID uuid.UUID) (*domain.Book, error)
		RetrievePage(ctx context.Context, filters *domain.BookFilters) (*domain.BookPage, error)
		Update(ctx context.Context, book *domain.Book) (*domain.Book, error)
		Remove(ctx context.Context, bookID uuid.UUID) error
	}
)
