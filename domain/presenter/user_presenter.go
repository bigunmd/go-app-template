package presenter

import "app/domain/model"

type UserPresenter interface {
	DefaultUser(u *model.User) *model.User
	DefaultUsers(us *model.Users) *model.Users
}
