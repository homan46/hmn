package service

import (
	"context"

	"codeberg.org/rchan/hmn/data/repository"
	"codeberg.org/rchan/hmn/model"
)

type UserService interface {
	AddUser(c context.Context, user *model.User) error
	GetUser(c context.Context, id int) (*model.User, error)
	GetAllUser(c context.Context) ([]*model.User, error)
	UpdateUser(c context.Context, user *model.User) error
	DeleteUser(c context.Context, id int) error

	CheckUserPassword(c context.Context, userName string, password string) (bool, error)
}

type UserServiceImpl struct {
	repo repository.RepositoryLayer
}

func NewUserService(repo repository.RepositoryLayer) UserService {
	return &UserServiceImpl{
		repo: repo,
	}
}

func (us *UserServiceImpl) AddUser(c context.Context, user *model.User) error {
	return us.repo.User().AddUser(c, user)
}
func (us *UserServiceImpl) GetUser(c context.Context, id int) (*model.User, error) {
	return us.repo.User().GetUser(c, id)
}

func (us *UserServiceImpl) GetAllUser(c context.Context) ([]*model.User, error) {
	return us.repo.User().GetAllUser(c)
}

func (us *UserServiceImpl) UpdateUser(c context.Context, user *model.User) error {
	return us.repo.User().UpdateUser(c, user)
}
func (us *UserServiceImpl) DeleteUser(c context.Context, id int) error {
	return us.repo.User().DeleteUser(c, id)
}

func (us *UserServiceImpl) CheckUserPassword(c context.Context, userName string, password string) (bool, error) {
	user, err := us.repo.User().GetUserByUserName(c, userName)
	if err != nil {
		return false, err
	}

	return user.CheckPassword(password), nil
}
