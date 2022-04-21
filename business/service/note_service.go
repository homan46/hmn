package service

import (
	"context"
	"errors"

	"codeberg.org/rchan/hmn/data/repository"
	"codeberg.org/rchan/hmn/helper"
	"codeberg.org/rchan/hmn/model"
)

var (
	ErrInvalidActionOnRoot = errors.New("cannot do action on root")
	ErrInvalidParent       = errors.New("parent is invalid")
	ErrInvalidIndex        = errors.New("index is invalid")
)

type NoteService interface {
	// only title,content,parentID and index in note parameter is used
	AddNote(c context.Context, note *model.Note) error
	GetNote(c context.Context, id int) (*model.Note, error)
	GetAllNote(c context.Context) ([]*model.Note, error)
	//will include the root note as well
	GetNoteUnder(c context.Context, rootID int) ([]*model.Note, error)
	UpdateNote(c context.Context, note *model.Note) error
	DeleteNote(c context.Context, id int) error

	MoveNote(c context.Context, id int, parentID int, index int) error
}

type NoteServiceImpl struct {
	repo repository.RepositoryLayer
}

func NewNoteService(repo repository.RepositoryLayer) NoteService {
	return &NoteServiceImpl{
		repo: repo,
	}
}

// only title,content,parentID and index in note parameter is used
func (ns *NoteServiceImpl) AddNote(c context.Context, note *model.Note) error {
	//validation
	if note.GetParentID() < 1 {
		return ErrInvalidParent
	}
	if note.GetIndex() < 0 {
		return ErrInvalidParent
	}

	userID, _, err := helper.ExtractContext(c)
	if err != nil {
		return err
	}

	note.Entity.SetUpdate(*userID)

	err = ns.repo.Note().AddNote(c, note)

	if err != nil {
		return err
	}

	return nil
}
func (ns *NoteServiceImpl) GetNote(c context.Context, id int) (*model.Note, error) {
	return ns.repo.Note().GetNote(c, id)
}
func (ns *NoteServiceImpl) GetAllNote(c context.Context) ([]*model.Note, error) {
	return ns.repo.Note().GetAllNote(c)
}

func (ns *NoteServiceImpl) GetNoteUnder(c context.Context, rootID int) ([]*model.Note, error) {
	return ns.repo.Note().GetNoteUnder(c, rootID)
}

func (ns *NoteServiceImpl) UpdateNote(c context.Context, note *model.Note) error {
	return ns.repo.Note().UpdateNote(c, note)
}
func (ns *NoteServiceImpl) DeleteNote(c context.Context, id int) error {
	if id == 1 {
		return ErrInvalidActionOnRoot
	}
	return ns.repo.Note().DeleteNote(c, id)
}
func (ns *NoteServiceImpl) MoveNote(c context.Context, id int, parentID int, index int) error {
	return errors.New("not implemented")
}
