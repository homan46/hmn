package helper

import (
	"context"
	"errors"

	"codeberg.org/rchan/hmn/constant"
	"github.com/jmoiron/sqlx"
)

func ExtractContext(c context.Context) (userID *int, tx *sqlx.Tx, err error) {
	if c.Value(constant.KeyofTx) == nil {
		tx = nil
	} else {
		tx1, ok := c.Value(constant.KeyofTx).(*sqlx.Tx)
		if !ok {
			return nil, nil, errors.New("extract tx fail")
		}
		tx = tx1
	}

	if c.Value(constant.KeyOfUserID) == nil {
		userID = nil
	} else {
		userID1, ok := c.Value(constant.KeyOfUserID).(int)
		if !ok {
			return nil, nil, errors.New("extract tx fail")
		}
		userID = &userID1
	}

	err = nil
	return
}
