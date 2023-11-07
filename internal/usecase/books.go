package usecase

import (
	"context"
	"fmt"
	"goapptemplate/internal/domain"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type booksUsecase struct {
	repo BooksRepo
	log  *logrus.Entry
}

// List implements Books.
func (u *booksUsecase) List(ctx context.Context, filters *domain.BookFilters) (*domain.BookPage, error) {
	err := filters.Validate()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", domain.ErrValidation, err)
	}
	p, err := u.repo.RetrievePage(ctx, filters)
	if err != nil {
		return nil, errors.Wrap(err, "cannot retrieve book page")
	}
	return p, nil
}

// Modify implements Books.
func (u *booksUsecase) Modify(ctx context.Context, book *domain.Book) (*domain.Book, error) {
	err := book.Validate()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", domain.ErrValidation, err)
	}
	b, err := u.repo.Update(ctx, book)
	if err != nil {
		return nil, errors.Wrap(err, "cannot update book")
	}
	return b, nil
}

// New implements Books.
func (u *booksUsecase) New(ctx context.Context, book *domain.Book) (*domain.Book, error) {
	err := book.Validate()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", domain.ErrValidation, err)
	}
	book.ID = uuid.New()
	b, err := u.repo.Store(ctx, book)
	if err != nil {
		return nil, errors.Wrap(err, "cannot store book")
	}
	return b, nil
}

// Remove implements Books.
func (u *booksUsecase) Remove(ctx context.Context, bookID uuid.UUID) error {
	err := u.repo.Remove(ctx, bookID)
	if err != nil {
		return errors.Wrapf(err, "cannot remove book with ID=%s", bookID)
	}
	return nil
}

// View implements Books.
func (u *booksUsecase) View(ctx context.Context, bookID uuid.UUID) (*domain.Book, error) {
	b, err := u.repo.Retrieve(ctx, bookID)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot retrieve book with ID=%s", bookID)
	}
	return b, nil
}

func NewBooks(repo BooksRepo, logger *logrus.Logger) Books {
	return &booksUsecase{
		repo: repo,
		log:  logger.WithField("layer", "internal.usecase.booksUsecase"),
	}
}
