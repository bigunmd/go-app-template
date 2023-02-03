package repository

import (
	"app/domain/model"

	"github.com/google/uuid"
)

type UserRepository interface {
	InsertUser(u *model.User) (*model.User, error)
	SelectUserById(id uuid.UUID) (*model.User, error)
	UpsertUser(u *model.User) (*model.User, error)
	DeleteUserById(id uuid.UUID) error
	SelectUser(*model.UserOptions) (*model.Users, error)
}
