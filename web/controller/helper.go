package controller

import (
	"errors"
	"math"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func readValueFromSession(c echo.Context, sessionName string, valueName string) (val interface{}, exists bool, err error) {
	sess, err := session.Get(sessionName, c)
	if err != nil {
		return nil, false, err
	}
	user_id, exists := sess.Values[valueName]
	return user_id, exists, nil
}

func readUserIDFromSession(c echo.Context) (userID int, err error) {
	val, exists, err := readValueFromSession(c, "credential", "user_id")
	if err != nil {
		return math.MaxInt, err
	}

	if !exists {
		return math.MaxInt, errors.New("not exist")
	}

	userID, ok := val.(int)

	if !ok {
		return math.MaxInt, errors.New("cannot parse userId in session to int")
	}

	return userID, nil

}
