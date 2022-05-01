package midd

import (
	"codeberg.org/rchan/hmn/business"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewAuth(bl business.BusinessLayer) echo.MiddlewareFunc {

	return middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		// Be careful to use constant time comparison to prevent timing attacks

		mycontext, tx, err := bl.GetContextForSystem()
		if err != nil {
			return false, err
		}

		pass, err := bl.User().CheckUserPassword(mycontext, username, password)
		if err != nil {
			tx.Rollback()
			return false, err
		}

		tx.Commit()

		return pass, nil
	})
}
