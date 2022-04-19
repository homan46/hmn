package business

import (
	"context"

	ser "codeberg.org/rchan/hmn/business/service"
	"codeberg.org/rchan/hmn/constant"
	repo "codeberg.org/rchan/hmn/data/repository"
	"github.com/jmoiron/sqlx"
)

type BusinessLayer interface {
	GetContextFor(userID int) (context.Context, *sqlx.Tx, error)
	Note() ser.NoteService
	User() ser.UserService
}

type BusinessLayerImpl struct {
	repoLayer   repo.RepositoryLayer
	noteService ser.NoteService
	userService ser.UserService
}

func NewBusunessLayer(db *sqlx.DB) BusinessLayer {
	repoLayer := repo.NewRepositoryLayer(db)
	return BusinessLayerImpl{
		repoLayer:   repoLayer,
		noteService: ser.NewNoteService(repoLayer),
		userService: ser.NewUserService(repoLayer),
	}
}

func (bl BusinessLayerImpl) GetContextFor(userID int) (context.Context, *sqlx.Tx, error) {
	c, tx, err := bl.repoLayer.GetContextWithTx()
	if err != nil {
		return nil, nil, err
	}
	return context.WithValue(
		c,
		constant.KeyOfUserID,
		userID,
	), tx, nil
}

func (bl BusinessLayerImpl) Note() ser.NoteService {
	return bl.noteService
}
func (bl BusinessLayerImpl) User() ser.UserService {
	return bl.userService
}
