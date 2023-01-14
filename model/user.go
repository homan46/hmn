package model

import (
	"encoding/base64"

	"codeberg.org/rchan/hmn/log"
	"go.uber.org/zap"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Entity
	userName string
	password string
}

func (u *User) GetUserName() string {
	return u.userName
}

func (u *User) SetUserName(name string) {
	u.userName = name
}

func (u *User) GetPassword() string {
	return u.password
}

func (u *User) SetPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.ZLog.Error("set password error(TODO: This log is not needed?)", zap.Error(err))
		return err
	}
	str := base64.StdEncoding.EncodeToString(bytes)
	u.password = str
	return nil
}

func (u *User) CheckPassword(password string) bool {
	bytes, err := base64.StdEncoding.DecodeString(u.password)
	if err != nil {
		log.ZLog.Error("check password error", zap.Error(err))
		return false
	}
	if bcrypt.CompareHashAndPassword(bytes, []byte(password)) != nil {
		return false
	}

	return true
}

type UserLikeRO interface {
	EntityLikeRO
	GetUserName() string
	GetPassword() string
}

func NewUserFrom(ul UserLikeRO) *User {
	newUser := User{}
	// mBy, mTime := ul.GetModifiedBy()
	// cBy, cTime := ul.GetCreatedBy()
	newUser.Entity = *NewEntity(ul.GetID(),
		ul.GetModifiedTime(), ul.GetModifiedBy(),
		ul.GetCreatedTime(), ul.GetCreatedBy())
	newUser.userName = ul.GetUserName()
	newUser.password = ul.GetPassword()

	return &newUser
}
