package user

import (
	"app/domain/model"
	"app/domain/presenter"
	"app/domain/repository"
	"app/domain/usecase"
	"app/pkg/logger"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type userUsecase struct {
	log            logger.Logger
	userRepository repository.UserRepository
	userPresenter  presenter.UserPresenter
}

// Create implements usecase.UserUsecase
func (uu *userUsecase) Create(u *model.UserBase) (*model.User, error) {
	usr := new(model.User)
	usr.UserBase = u
	usr, err := uu.userRepository.InsertUser(usr)
	if err != nil {
		return nil, err
	}
	return uu.userPresenter.DefaultUser(usr), nil
}

// DeleteById implements usecase.UserUsecase
func (uu *userUsecase) DeleteById(id uuid.UUID) error {
	err := uu.userRepository.DeleteUserById(id)
	if err != nil {
		return err
	}
	return nil
}

// Get implements usecase.UserUsecase
func (uu *userUsecase) Get(opts *model.UserOptions) (*model.Users, error) {
	us, err := uu.userRepository.SelectUser(opts)
	if err != nil {
		return nil, err
	}
	return uu.userPresenter.DefaultUsers(us), nil
}

// GetById implements usecase.UserUsecase
func (uu *userUsecase) GetById(id uuid.UUID) (*model.User, error) {
	u, err := uu.userRepository.SelectUserById(id)
	if err != nil {
		return nil, err
	}
	return uu.userPresenter.DefaultUser(u), nil
}

// Update implements usecase.UserUsecase
func (uu *userUsecase) Update(u *model.User) (*model.User, error) {
	u, err := uu.userRepository.UpsertUser(u)
	if err != nil {
		return nil, err
	}
	return uu.userPresenter.DefaultUser(u), err
}

func NewUserUsecase(db *sqlx.DB, logger logger.Logger) usecase.UserUsecase {
	ur := NewUserRepository(db, logger)
	up := NewUserPresenter(logger)
	return &userUsecase{
		log:            logger,
		userRepository: ur,
		userPresenter:  up,
	}
}
