package service

import (
	"context"
	"errors"

	"codeberg.org/rchan/hmn/model"
	"codeberg.org/rchan/hmn/repository"
)

var (
	ErrUserNotExist = errors.New("user does not exist")
)

type UserService interface {
	AddUser(c context.Context, user *model.User) error
	GetUser(c context.Context, id int) (*model.User, error)
	GetUserByUserName(c context.Context, userName string) (*model.User, error)
	GetAllUser(c context.Context) ([]*model.User, error)
	UpdateUser(c context.Context, user *model.User) error
	DeleteUser(c context.Context, id int) error

	CheckUserPassword(c context.Context, userName string, password string) (bool, error)
	SetUserPassword(c context.Context, userName string, password string) error
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

func (us *UserServiceImpl) GetUserByUserName(c context.Context, userName string) (*model.User, error) {
	return us.repo.User().GetUserByUserName(c, userName)
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
		if err.Error() == "sql: no rows in result set" {
			return false, ErrUserNotExist
		}
		return false, err
	}

	return user.CheckPassword(password), nil
}

func (us *UserServiceImpl) SetUserPassword(c context.Context, userName string, password string) error {
	user, err := us.repo.User().GetUserByUserName(c, userName)
	if err != nil {
		return err
	}

	user.SetPassword(password)
	err = us.repo.User().UpdateUser(c, user)
	if err != nil {
		return err
	}

	return nil
}
