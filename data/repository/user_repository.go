package repository

import (
	"context"
	"errors"
	"time"

	"codeberg.org/rchan/hmn/helper"
	"codeberg.org/rchan/hmn/model"
	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	GetUser(c context.Context, id int) (*model.User, error)
	GetAllUser(c context.Context) ([]*model.User, error)
	AddUser(c context.Context, user *model.User) error
	UpdateUser(c context.Context, user *model.User) error
	DeleteUser(c context.Context, id int) error
}

type SqlxUserRepository struct {
	db *sqlx.DB
}

func (r SqlxUserRepository) GetUser(c context.Context, id int) (*model.User, error) {
	_, tx, err := helper.ExtractContext(c)
	if err != nil {
		return nil, errors.New("get Context data fail")
	}

	//user := model.User{}
	userEntity := UserEntity{}
	err = tx.Get(&userEntity, "select * from user where id = ?", id)
	if err != nil {
		return nil, err
	}

	user := model.NewUserFrom(&userEntity)

	return user, nil
}

func (r SqlxUserRepository) GetAllUser(c context.Context) ([]*model.User, error) {
	_, tx, err := helper.ExtractContext(c)
	if err != nil {
		return nil, errors.New("get Context data fail")
	}

	//noteEntities := model.Note{}
	userEntities := []UserEntity{}

	err = tx.Select(&userEntities, `
	select * from note
	`)

	if err != nil {
		return nil, err
	}

	users := make([]*model.User, len(userEntities))
	for _, e := range userEntities {
		users = append(users, model.NewUserFrom(&e))
	}

	return users, nil
}

func (r SqlxUserRepository) AddUser(c context.Context, user *model.User) error {
	userID, tx, err := helper.ExtractContext(c)
	if err != nil {
		return errors.New("get Context data fail")
	}

	result, err := tx.Exec(`
	insert into user (
		created_time,created_by,
		modified_time,modified_by,
		user_name, password
	) values (
		datetime(),?,
		datetime(),?,
		?,?
	)
	`, userID, userID, user.GetUserName(), user.GetPassword())

	if err != nil {
		return err
	}
	//TODO: should load the model from database
	newID, err := result.LastInsertId()
	user.SetID(int(newID)) //TODO cast without check
	return nil
}
func (r SqlxUserRepository) UpdateUser(c context.Context, user *model.User) error {
	userID, tx, err := helper.ExtractContext(c)
	if err != nil {
		return errors.New("get Context data fail")
	}

	_, err = tx.Exec(`
	update user set modified_time = datetime(), modified_by = ?,
	username = ?, password = ?
	where id = ?
	`, userID, user.GetUserName(), user.GetPassword(), user.GetID())

	if err != nil {
		return err
	}

	return nil

}
func (r SqlxUserRepository) DeleteUser(c context.Context, id int) error {
	_, tx, err := helper.ExtractContext(c)
	if err != nil {
		return errors.New("get Context data fail")
	}
	_, err = tx.Exec(`
	delete from user
	where id = ?
	`, id)

	if err != nil {
		return err
	}

	return nil
}

type UserEntity struct {
	ID           int        `db:"id"`
	ModifiedTime *time.Time `db:"modified_time"`
	ModifiedBy   int        `db:"modified_by"`
	CreatedTime  *time.Time `db:"created_time"`
	CreatedBy    int        `db:"created_by"`

	UserName string `db:"user_name"`
	Password string `db:"password"`
}

func (e *UserEntity) GetID() int {
	return e.ID
}
func (e *UserEntity) GetCreatedBy() (userID int, createdTime time.Time) {
	return e.CreatedBy, *e.CreatedTime
}
func (e *UserEntity) GetModifiedBy() (userID int, modifiedTime time.Time) {
	return e.ModifiedBy, *e.ModifiedTime
}

func (n *UserEntity) GetUserName() string {
	return n.UserName
}
func (n *UserEntity) GetPassword() string {
	return n.Password
}
