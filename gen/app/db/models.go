// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0

package db

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Book struct {
	ID          pgtype.UUID
	Name        string
	Description pgtype.Text
	CreatedAt   pgtype.Timestamptz
	UpdatedAt   pgtype.Timestamptz
}