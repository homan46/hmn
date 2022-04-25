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

	//check new parent exist
	_, err = ns.repo.Note().GetNote(c, parentID)
	if err != nil {
		return err
	}

	if note.GetParentID() != parentID {
		//TODO: check move under itself
		thisAndUnder, err := ns.repo.Note().GetNoteUnder(c, id)

		for _, child := range thisAndUnder[1:] {
			if parentID == child.GetID() {
				return errors.New("new parent is child")
			}
		}

		//index -1 to all the old sibling that is after current note
		updateWaitList := make([]*model.Note, 0)

		oldSibling, err := ns.GetNoteUnder(c, note.GetParentID())
		if err != nil {
			return err
		}
		sort.Slice(oldSibling[1:], func(i, j int) bool {
			return oldSibling[i+1].GetIndex() < oldSibling[j+1].GetIndex()
		})

		for _, sib := range oldSibling[1:] {
			if sib.GetIndex() > note.GetIndex() {
				sib.SetIndex(sib.GetIndex() - 1)
				updateWaitList = append(updateWaitList, sib)
			}
		}

		//index + 1 to all the new sibling that is after current note
		newSibling, err := ns.GetNoteUnder(c, parentID)
		if err != nil {
			return err
		}
		sort.Slice(newSibling[1:], func(i, j int) bool {
			return newSibling[i+1].GetIndex() < newSibling[j+1].GetIndex()
		})

		for _, sib := range newSibling[1:] {
			if sib.GetIndex() > index {
				sib.SetIndex(sib.GetIndex() + 1)
				updateWaitList = append(updateWaitList, sib)
			}
		}

		note.SetIndex(index)
		note.SetParentID(parentID)

		//ns.repo.Note().UpdateNote(c, note)
		updateWaitList = append(updateWaitList, note)

		for _, n := range updateWaitList {
			err = ns.repo.Note().UpdateNote(c, n)
			if err != nil {
				return err
			}
		}
		return nil

	} else {
		updateWaitList := make([]*model.Note, 0)

		oldIdx := note.GetIndex()
		newIdx := index

		sibling, err := ns.GetNoteUnder(c, note.GetParentID())
		if err != nil {
			return err
		}
		sort.Slice(sibling[1:], func(i, j int) bool {
			return sibling[i+1].GetIndex() < sibling[j+1].GetIndex()
		})

		if newIdx < oldIdx {
			//move front,
			//all between two Index  + 1

			for i, n := range sibling[1:] {
				if i >= newIdx && i < oldIdx {
					n.SetIndex(n.GetIndex() + 1)
					updateWaitList = append(updateWaitList, n)
				}
			}

			note.SetIndex(index)
			note.SetParentID(parentID)
			updateWaitList = append(updateWaitList, note)

		} else if newIdx > oldIdx {
			//move back,
			//all between two Index  - 1
			for i, n := range sibling[1:] {
				if i >= newIdx && i < oldIdx {
					n.SetIndex(n.GetIndex() - 1)
					updateWaitList = append(updateWaitList, n)
				}
			}

			note.SetIndex(index)
			note.SetParentID(parentID)
			updateWaitList = append(updateWaitList, note)

		} else {
			return nil
		}

		for _, n := range updateWaitList {
			err = ns.repo.Note().UpdateNote(c, n)
			if err != nil {
				return err
			}
		}

		return nil

	}

}
