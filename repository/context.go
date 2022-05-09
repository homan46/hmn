package repository

import (
	"context"

	"codeberg.org/rchan/hmn/constant"
	"github.com/jmoiron/sqlx"
)

type RepositoryLayer interface {
	GetContextWithTx() (context.Context, *sqlx.Tx, error)
	Note() NoteRepository
	User() UserRepository
}

type SqlxRepositoryLayer struct {
	noteRepo NoteRepository
	userRepo UserRepository
	db       *sqlx.DB
}

func (l SqlxRepositoryLayer) GetContextWithTx() (context.Context, *sqlx.Tx, error) {

	tx, err := l.db.Beginx()
	if err != nil {
		return nil, nil, err
	}

	return context.WithValue(
		context.Background(),
		constant.KeyofTx, tx,
	), tx, nil
}

func (l SqlxRepositoryLayer) Note() NoteRepository {
	return l.noteRepo
}
func (l SqlxRepositoryLayer) User() UserRepository {
	return l.userRepo
}

func NewRepositoryLayer(db *sqlx.DB) RepositoryLayer {
	repoLayer := SqlxRepositoryLayer{
		noteRepo: SqlxNoteRepository{db},
		userRepo: SqlxUserRepository{db},
		db:       db,
	}

	return repoLayer
}
