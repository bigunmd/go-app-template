package repo

import (
	"context"
	"goapptemplate/gen/app/db"
	"goapptemplate/internal/domain"
	"goapptemplate/internal/usecase"
	"goapptemplate/pkg/postgres"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type booksPostgresRepo struct {
	postgres.DB
	log *logrus.Entry
}

// Remove implements usecase.BooksRepo.
func (repo *booksPostgresRepo) Remove(ctx context.Context, bookID uuid.UUID) error {
	conn, tx, err := repo.BeginTx(ctx)
	if err != nil {
		return errors.Wrap(err, "cannot begin tx")
	}
	defer conn.Release()
	defer tx.Rollback(ctx)

	q := db.New(conn).WithTx(tx)

	err = q.DeleteBookWhereID(ctx, pgtype.UUID{
		Bytes: bookID,
		Valid: true,
	})
	if err != nil {
		return errors.Wrapf(err, "cannot delete book where ID=%s", bookID)
	}

	err = repo.EndTx(ctx, tx)
	if err != nil {
		return errors.Wrap(err, "cannot end tx")
	}

	return nil
}

// Retrieve implements usecase.BooksRepo.
func (repo *booksPostgresRepo) Retrieve(ctx context.Context, bookID uuid.UUID) (*domain.Book, error) {
	conn, tx, err := repo.BeginTx(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "cannot begin tx")
	}
	defer conn.Release()
	defer tx.Rollback(ctx)

	q := db.New(conn).WithTx(tx)

	row, err := q.SelectBookWhereID(ctx, pgtype.UUID{
		Bytes: bookID,
		Valid: true,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrBookNotFound
		}
		return nil, errors.Wrapf(err, "cannot select book where ID=%s", bookID)
	}
	book := &domain.Book{
		ID:          row.ID.Bytes,
		Name:        row.Name,
		Description: row.Description.String,
		CreatedAt:   row.CreatedAt.Time,
		UpdatedAt:   row.UpdatedAt.Time,
	}

	err = repo.EndTx(ctx, tx)
	if err != nil {
		return nil, errors.Wrap(err, "cannot end tx")
	}

	return book, nil
}

// RetrievePage implements usecase.BooksRepo.
func (repo *booksPostgresRepo) RetrievePage(ctx context.Context, filters *domain.BookFilters) (*domain.BookPage, error) {
	conn, tx, err := repo.BeginTx(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "cannot begin tx")
	}
	defer conn.Release()
	defer tx.Rollback(ctx)

	q := db.New(conn).WithTx(tx)

	total, err := q.SelectBooksCount(ctx, db.SelectBooksCountParams{
		Name: "%" + filters.Name + "%",
		Description: pgtype.Text{
			String: "%" + filters.Description + "%",
			Valid:  true,
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "cannot select books count")
	}
	rows, err := q.SelectBooks(ctx, db.SelectBooksParams{
		Name: "%" + filters.Name + "%",
		Description: pgtype.Text{
			String: "%" + filters.Description + "%",
			Valid:  true,
		},
		Ofst: filters.Offset,
		Lim:  filters.Limit,
	})
	if err != nil {
		return nil, errors.Wrap(err, "cannot select books")
	}
	var books []*domain.Book
	for _, b := range rows {
		book := &domain.Book{
			ID:          b.ID.Bytes,
			Name:        b.Name,
			Description: b.Description.String,
			CreatedAt:   b.CreatedAt.Time,
			UpdatedAt:   b.UpdatedAt.Time,
		}
		books = append(books, book)
	}

	err = repo.EndTx(ctx, tx)
	if err != nil {
		return nil, errors.Wrap(err, "cannot end tx")
	}

	return &domain.BookPage{
		Page: domain.Page{
			Total:  total,
			Limit:  filters.Limit,
			Offset: filters.Offset,
			Metadata: map[string]interface{}{
				"description": "filtered page of books",
			},
		},
		Data: books,
	}, nil
}

// Store implements usecase.BooksRepo.
func (repo *booksPostgresRepo) Store(ctx context.Context, book *domain.Book) (*domain.Book, error) {
	conn, tx, err := repo.BeginTx(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "cannot begin tx")
	}
	defer conn.Release()
	defer tx.Rollback(ctx)

	q := db.New(conn).WithTx(tx)

	row, err := q.InsertBook(ctx, db.InsertBookParams{
		ID: pgtype.UUID{
			Bytes: book.ID,
			Valid: true,
		},
		Name: book.Name,
		Description: pgtype.Text{
			String: book.Description,
			Valid:  true,
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "cannot insert book")
	}
	book = &domain.Book{
		ID:          row.ID.Bytes,
		Name:        row.Name,
		Description: row.Description.String,
		CreatedAt:   row.CreatedAt.Time,
		UpdatedAt:   row.UpdatedAt.Time,
	}

	err = repo.EndTx(ctx, tx)
	if err != nil {
		return nil, errors.Wrap(err, "cannot end tx")
	}

	return book, nil
}

// Update implements usecase.BooksRepo.
func (repo *booksPostgresRepo) Update(ctx context.Context, book *domain.Book) (*domain.Book, error) {
	conn, tx, err := repo.BeginTx(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "cannot begin tx")
	}
	defer conn.Release()
	defer tx.Rollback(ctx)

	q := db.New(conn).WithTx(tx)

	err = q.UpdateBookWhereID(ctx, db.UpdateBookWhereIDParams{
		Name: book.Name,
		Description: pgtype.Text{
			String: book.Description,
			Valid:  true,
		},
		ID: pgtype.UUID{
			Bytes: book.ID,
			Valid: true,
		},
	})
	if err != nil {
		return nil, errors.Wrapf(err, "cannot update book where ID=%s", book.ID)
	}
	row, err := q.SelectBookWhereID(ctx, pgtype.UUID{
		Bytes: book.ID,
		Valid: true,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrBookNotFound
		}
		return nil, errors.Wrapf(err, "cannot select book where ID=%s", book.ID)
	}
	book = &domain.Book{
		ID:          row.ID.Bytes,
		Name:        row.Name,
		Description: row.Description.String,
		CreatedAt:   row.CreatedAt.Time,
		UpdatedAt:   row.UpdatedAt.Time,
	}

	err = repo.EndTx(ctx, tx)
	if err != nil {
		return nil, errors.Wrap(err, "cannot end tx")
	}

	return book, nil
}

func NewBooksPostgresRepo(db postgres.DB, logger *logrus.Logger) usecase.BooksRepo {
	return &booksPostgresRepo{
		DB:  db,
		log: logger.WithField("layer", "internal.usecase.repo.booksPostgresRepo"),
	}
}
