package usecase

import (
	"app/domain/model"

	"github.com/google/uuid"
)

type UserUsecase interface {
	Create(u *model.UserBase) (*model.User, error)
	GetById(id uuid.UUID) (*model.User, error)
	Update(u *model.User) (*model.User, error)
	DeleteById(id uuid.UUID) error
	Get(opts *model.UserOptions) (*model.Users, error)
}
