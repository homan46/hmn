package model

import "time"

type EntityEmbed struct {
	id           int
	modifiedTime time.Time
	modifiedBy   int
	createdTime  time.Time
	createdBy    int
}

func (e *EntityEmbed) GetID() int {
	return e.id
}
func (e *EntityEmbed) SetID(id int) {
	e.id = id
}

func (e *EntityEmbed) GetCreatedBy() (userID int, createdTime time.Time) {
	return e.createdBy, e.createdTime
}

func (e *EntityEmbed) GetModifiedBy() (userID int, modifiedTime time.Time) {
	return e.modifiedBy, e.modifiedTime
}

func (e *EntityEmbed) SetUpdate(userID int) {
	e.modifiedTime = time.Now().UTC()
	e.modifiedBy = userID
}

func NewEntityEmbed(
	id int,
	modifiedTime time.Time, modifiedBy int,
	createdTime time.Time, createdBy int) *EntityEmbed {
	return &EntityEmbed{
		id:           id,
		modifiedTime: modifiedTime,
		modifiedBy:   modifiedBy,
		createdTime:  createdTime,
		createdBy:    createdBy,
	}
}

type EntityEmbedLikeRO interface {
	GetID() int
	//SetID(id int)
	GetCreatedBy() (userID int, createdTime time.Time)
	GetModifiedBy() (userID int, modifiedTime time.Time)
	//SetUpdate(userID int)
}
