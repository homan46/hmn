package dto

import "time"

type entityIDField struct {
	ID int `db:"id" json:"id"`
}

func (e *entityIDField) GetID() int {
	return e.ID
}

func (e *entityIDField) SetID(id int) {
	e.ID = id
}

//////

type entityModifiedTimeField struct {
	ModifiedTime time.Time `db:"modified_time" json:"modifiedTime"`
}

func (e *entityModifiedTimeField) GetModifiedTime() time.Time {
	return e.ModifiedTime
}

func (e *entityModifiedTimeField) SetModifiedTime(modifiedTime time.Time) {
	e.ModifiedTime = modifiedTime
}

/////
type entityModifiedByField struct {
	ModifiedBy int `db:"modified_by" json:"modifiedBy"`
}

func (e *entityModifiedByField) GetModifiedBy() int {
	return e.ModifiedBy
}

func (e *entityModifiedByField) SetModifiedBy(modifiedBy int) {
	e.ModifiedBy = modifiedBy
}

///////////////////
//////////
////////////
//////////////

type entityCreatedTimeField struct {
	CreatedTime time.Time `db:"created_time" json:"createdTime"`
}

func (e *entityCreatedTimeField) GetCreatedTime() time.Time {
	return e.CreatedTime
}

func (e *entityCreatedTimeField) SetCreatedTime(createdTime time.Time) {
	e.CreatedTime = createdTime
}

/////
type entityCreatedByField struct {
	CreatedBy int `db:"created_by" json:"createdBy"`
}

func (e *entityCreatedByField) GetCreatedBy() int {
	return e.CreatedBy
}

func (e *entityCreatedByField) SetCreatedBy(createdBy int) {
	e.CreatedBy = createdBy
}
