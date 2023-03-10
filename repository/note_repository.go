package repository

import (
	"context"
	"errors"

	"codeberg.org/rchan/hmn/dto"
	"codeberg.org/rchan/hmn/helper"
	"codeberg.org/rchan/hmn/log"
	"codeberg.org/rchan/hmn/model"
	"github.com/jmoiron/sqlx"
)

type NoteRepository interface {
	GetNote(c context.Context, id int) (*model.Note, error)
	GetAllNote(c context.Context) ([]*model.Note, error)
	AddNote(c context.Context, note *model.Note) error
	UpdateNote(c context.Context, note *model.Note) error
	DeleteNote(c context.Context, id int) error

	// return root note and all note under it
	// the result is sorted by depth so root note is always the first
	GetNoteUnder(c context.Context, rootNoteID int, depth int) ([]*model.Note, error)
}

type SqlxNoteRepository struct {
	db *sqlx.DB
}

func (r SqlxNoteRepository) GetNote(c context.Context, id int) (*model.Note, error) {
	_, tx, err := helper.ExtractContext(c)
	if err != nil {
		return nil, errors.New("get Context data fail")
	}

	noteEntity := dto.NoteEntityDto{}

	err = tx.Get(&noteEntity, `
	select *
	from note 
	where id = ?`, id)
	if err != nil {
		return nil, err
	}

	note := model.NewNoteFrom(&noteEntity)

	return note, nil
}

// GetAllNote will return all note as a list.
// the result is not ordered.
func (r SqlxNoteRepository) GetAllNote(c context.Context) ([]*model.Note, error) {
	log.ZLog.Debug("SqlxNoteRepository:GetAllNote")
	_, tx, err := helper.ExtractContext(c)
	if err != nil {
		return nil, errors.New("get Context data fail")
	}

	//noteEntities := model.Note{}
	noteEntities := []dto.NoteEntityDto{}

	err = tx.Select(&noteEntities, `
	select *
	from note `)
	if err != nil {
		return nil, err
	}

	notes := make([]*model.Note, 0)
	for _, e := range noteEntities {
		notes = append(notes, model.NewNoteFrom(&e))
	}

	return notes, nil
}

func (r SqlxNoteRepository) AddNote(c context.Context, note *model.Note) error {
	log.ZLog.Sugar().Debugf("SqlxNoteRepository:AddNote %v\n", note)

	userID, tx, err := helper.ExtractContext(c)
	if err != nil {
		return errors.New("get Context data fail")
	}

	result, err := tx.Exec(`
	insert into note (created_time,created_by,
		modified_time,modified_by,
		title,content,parent_id,idx)
	values (
		datetime(),?,
		datetime(),?,
		?,?,?,?
	)`, userID, userID, note.GetTitle(), note.GetContent(), note.GetParentID(), note.GetIndex())

	if err != nil {
		return err
	}

	//TODO: should load the model from database
	newID, err := result.LastInsertId()
	note.SetID(int(newID)) //TODO cast without check
	return nil
}
func (r SqlxNoteRepository) UpdateNote(c context.Context, note *model.Note) error {
	log.ZLog.Sugar().Debugf("SqlxNoteRepository:UpdateNote %v\n", note)

	userID, tx, err := helper.ExtractContext(c)
	if err != nil {
		return errors.New("get Context data fail")
	}

	_, err = tx.Exec(`
	update note set modified_time = datetime(), modified_by = ?,
	title = ?, content = ?, parent_id = ?,idx = ?
	where id = ?
	`, userID, note.GetTitle(), note.GetContent(), note.GetParentID(), note.GetIndex(), note.GetID())

	if err != nil {
		return err
	}

	return nil
}
func (r SqlxNoteRepository) DeleteNote(c context.Context, id int) error {
	log.ZLog.Debug("SqlxNoteRepository:DeleteNote")
	_, tx, err := helper.ExtractContext(c)
	if err != nil {
		return errors.New("get Context data fail")
	}

	_, err = tx.Exec(`
	delete from note
	where id = ?
	`, id)

	if err != nil {
		return err
	}

	return nil
}

func (r SqlxNoteRepository) GetNoteUnder(c context.Context, rootNoteID int, depth int) ([]*model.Note, error) {
	log.ZLog.Debug("SqlxNoteRepository:GetNoteUnder")
	_, tx, err := helper.ExtractContext(c)
	if err != nil {
		return nil, errors.New("get Context data fail")
	}

	//noteEntities := model.Note{}
	noteEntities := []dto.NoteEntityDto{}

	err = tx.Select(&noteEntities, `	
		WITH RECURSIVE
		note_tree (
			depth,
			id ,
			created_time ,
			created_by	,
			modified_time ,
			modified_by ,
		
			title ,
			content ,
			parent_id ,
			idx 
		) as (
			select 0 as depth, * from note where id = ?
			union
			select depth+1 ,n.* from note n join note_tree t 
			on n.parent_id  = t.id
		)
		select  
			id ,
			created_time ,
			created_by	,
			modified_time ,
			modified_by ,
		
			title ,
			content ,
			parent_id ,
			idx 
		from note_tree 
		where depth <= ?
		order by depth,parent_id,idx;
	`, rootNoteID, depth) //TODO: add ordering for each layer

	if err != nil {
		return nil, err
	}

	notes := make([]*model.Note, 0)
	for _, e := range noteEntities {
		notes = append(notes, model.NewNoteFrom(&e))
	}

	return notes, nil
}
