package domain

import (
	"time"

	"github.com/google/uuid"
)

type Book struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (b Book) Validate() error {
	if b.Name == "" ||
		len(b.Name) > 255 {
		return ErrBookName
	}
	return nil
}

type BookFilters struct {
	Filters
	Name        string `json:"name" query:"name"`
	Description string `json:"description" query:"description"`
}

func (f *BookFilters) Validate() error {
	return f.Filters.Validate()
}

type BookPage struct {
	Page
	Data []*Book `json:"data"`
}
