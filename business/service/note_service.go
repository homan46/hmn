package service

import (
	"context"
	"errors"
	"math"

	"codeberg.org/rchan/hmn/constant"
	"codeberg.org/rchan/hmn/data/repository"
	"codeberg.org/rchan/hmn/helper"
	"codeberg.org/rchan/hmn/model"
)

var (
	ErrInvalidActionOnRoot = errors.New("cannot do action on root")
	ErrInvalidParent       = errors.New("parent is invalid")
	ErrInvalidIndex        = errors.New("index is invalid")
	ErrFieldType           = errors.New("field type mismatch")
)

type NoteService interface {
	// only title,content,parentID and index in note parameter is used
	AddNote(c context.Context, note *model.Note) error
	GetNote(c context.Context, id int) (*model.Note, error)
	GetAllNote(c context.Context) ([]*model.Note, error)
	//will include the root note as well
	GetNoteUnder(c context.Context, rootID int) ([]*model.Note, error)
	UpdateNote(c context.Context, note *model.Note) error
	PatchNote(c context.Context, id int, input map[string]interface{}) error
	DeleteNote(c context.Context, id int) error

	//MoveNote(c context.Context, id int, parentID int, index int) error
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

	//as long as parentID is valid, action should not fail
	targetParent, err := ns.repo.Note().GetNote(c, note.GetParentID())
	if err != nil {
		return ErrInvalidParent
	}

	//set the Index field of the to-be-added note
	currentChildren, err := ns.repo.Note().GetNoteUnder(c, targetParent.GetID(), 1)
	directChildCount := 0
	for _, child := range currentChildren[1:] {
		if child.GetParentID() == targetParent.GetID() {
			directChildCount += 1
		}
		//TODO: stop early
	}

	note.SetIndex(directChildCount)

	//actual work
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
	return ns.repo.Note().GetNoteUnder(c, rootID, math.MaxInt)
}

func (ns *NoteServiceImpl) UpdateNote(c context.Context, note *model.Note) error {

	oldValue, err := ns.repo.Note().GetNote(c, note.GetID())
	if err != nil {
		return err
	}

	//TODO: should  not move under its children

	//update title and content first
	oldValue.SetTitle(note.GetTitle())
	oldValue.SetContent(note.GetContent())

	err = ns.repo.Note().UpdateNote(c, oldValue)
	if err != nil {
		return err
	}

	//update position
	if oldValue.GetParentID() != note.GetParentID() || oldValue.GetIndex() != note.GetIndex() {
		//do nothing
	} else {
		err = ns.moveNote(c, oldValue.GetID(), note.GetParentID(), note.GetIndex())
		if err != nil {
			return err
		}

	}

	return nil

}

func (ns *NoteServiceImpl) PatchNote(c context.Context, id int, input map[string]interface{}) error {

	oldValue, err := ns.repo.Note().GetNote(c, id)
	newValue := *oldValue

	newValueWithoutPositionChange := *oldValue

	if err != nil {
		return err
	}

	//TODO: should  not move under its children

	positionChange := false

	for field, value := range input {
		switch field {
		case "title":
			v, ok := value.(string)
			if !ok {
				return ErrFieldType
			}
			newValue.SetTitle(v)
			newValueWithoutPositionChange.SetTitle(v)
			break
		case "content":
			v, ok := value.(string)
			if !ok {
				return ErrFieldType
			}
			newValue.SetContent(v)
			newValueWithoutPositionChange.SetContent(v)
			break
		case "parentId":

			f, ok := value.(float64)

			if !ok {
				return ErrFieldType
			}
			v := int(f)

			if oldValue.GetParentID() != v {
				positionChange = true
			}
			newValue.SetParentID(v)
			break
		case "index":
			f, ok := value.(float64)

			if !ok {
				return ErrFieldType
			}

			v := int(f)
			if oldValue.GetIndex() != v {
				positionChange = true
			}
			newValue.SetIndex(v)
			break
		}
	}

	err = ns.repo.Note().UpdateNote(c, &newValueWithoutPositionChange)
	if err != nil {
		return err
	}

	if !positionChange {
		return nil
	}

	//update position

	err = ns.moveNote(c, oldValue.GetID(), newValue.GetParentID(), newValue.GetIndex())
	if err != nil {
		return err
	}

	return nil
}

func (ns *NoteServiceImpl) DeleteNote(c context.Context, id int) error {
	if id == constant.RootNoteID {
		return ErrInvalidActionOnRoot
	}
	return ns.repo.Note().DeleteNote(c, id)
}

func (ns *NoteServiceImpl) moveNote(c context.Context, id int, parentID int, index int) error {
	note, err := ns.repo.Note().GetNote(c, id)
	if err != nil {
		return err
	}

	pendingUpdateList := make([]*model.Note, 0)

	if note.GetParentID() != parentID {
		//check new parent is not child
		thisAndUnder, err := ns.repo.Note().GetNoteUnder(c, id, math.MaxInt)

		for _, child := range thisAndUnder[1:] {
			if parentID == child.GetID() {
				return errors.New("new parent is child")
			}
		}

		//remove from old position, change index of old sibling
		oldPosition, err := ns.repo.Note().GetNoteUnder(c, note.GetParentID(), 1)
		if err != nil {
			return err
		}

		oldSibling := oldPosition[1:]
		for _, sibling := range oldSibling {
			if sibling.GetIndex() > note.GetIndex() {
				sibling.SetIndex(sibling.GetIndex() - 1)
				pendingUpdateList = append(pendingUpdateList, sibling)
			}
		}
		//insert to new position position, change index of new sibling
		newPosition, err := ns.repo.Note().GetNoteUnder(c, parentID, 1)
		if err != nil {
			return err
		}

		newSibling := newPosition[1:]
		for _, sibling := range newSibling {
			if sibling.GetIndex() >= note.GetIndex() {
				sibling.SetIndex(sibling.GetIndex() + 1)
				pendingUpdateList = append(pendingUpdateList, sibling)
			}
		}
	} else {
		if note.GetIndex() == index {
			return nil
		}

		oldPosition, err := ns.repo.Note().GetNoteUnder(c, note.GetParentID(), 1)
		if err != nil {
			return err
		}

		oldSibling := oldPosition[1:]
		for _, sibling := range oldSibling {
			if note.GetIndex() > index { //move to tail
				if sibling.GetIndex() <= index && sibling.GetIndex() > note.GetIndex() {
					sibling.SetIndex(sibling.GetIndex() - 1)
					pendingUpdateList = append(pendingUpdateList, sibling)
				}
			} else if note.GetIndex() < index { //move to 0
				if sibling.GetIndex() >= index && sibling.GetIndex() < note.GetIndex() {
					sibling.SetIndex(sibling.GetIndex() + 1)
					pendingUpdateList = append(pendingUpdateList, sibling)
				}
			}
		}
	}

	note.SetParentID(parentID)
	note.SetIndex(index)
	pendingUpdateList = append(pendingUpdateList, note)

	for _, n := range pendingUpdateList {
		err = ns.repo.Note().UpdateNote(c, n)
		if err != nil {
			return err
		}
	}

	return nil

}
