package user

import (
	"app/domain/model"
	"app/domain/presenter"
	"app/pkg/logger"
)

type userPresenter struct {
	log logger.Logger
}

// DefaultUser implements presenter.UserPresenter
func (*userPresenter) DefaultUser(u *model.User) *model.User {
	return u
}

// DefaultUsers implements presenter.UserPresenter
func (*userPresenter) DefaultUsers(us *model.Users) *model.Users {
	return us
}

func NewUserPresenter(logger logger.Logger) presenter.UserPresenter {
	return &userPresenter{
		log: logger,
	}
}
