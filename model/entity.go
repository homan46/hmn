package model

import "time"

type Entity struct {
	id           int
	modifiedTime time.Time
	modifiedBy   int
	createdTime  time.Time
	createdBy    int
}

func (e *Entity) GetID() int {
	return e.id
}
func (e *Entity) SetID(id int) {
	e.id = id
}

func (e *Entity) GetCreatedBy() (userID int) {
	return e.createdBy
}

func (e *Entity) GetCreatedTime() (createdTime time.Time) {
	return e.createdTime
}

func (e *Entity) GetModifiedBy() (userID int) {
	return e.modifiedBy
}

func (e *Entity) GetModifiedTime() (modifiedTime time.Time) {
	return e.modifiedTime
}

func (e *Entity) SetUpdate(userID int) {
	e.modifiedTime = time.Now().UTC()
	e.modifiedBy = userID
}

func NewEntity(
	id int,
	modifiedTime time.Time, modifiedBy int,
	createdTime time.Time, createdBy int) *Entity {
	return &Entity{
		id:           id,
		modifiedTime: modifiedTime,
		modifiedBy:   modifiedBy,
		createdTime:  createdTime,
		createdBy:    createdBy,
	}
}

type EntityLikeRO interface {
	GetID() int
	GetCreatedBy() (userID int)
	GetCreatedTime() (createdTime time.Time)
	GetModifiedBy() (userID int)
	GetModifiedTime() (modifiedTime time.Time)
	//SetUpdate(userID int)
}

/*
type EntityEmbedLikeWO interface {
	SetID(id int)
	SetCreatedBy(userID int, createdTime time.Time)
	SetModifiedBy(userID int, modifiedTime time.Time)
}
*/
