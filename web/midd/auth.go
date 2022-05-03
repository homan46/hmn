package midd

import (
	"codeberg.org/rchan/hmn/business"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
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

		if pass {

			user, err := bl.User().GetUserByUserName(mycontext, username)
			if err != nil {
				tx.Rollback()
				return false, err
			}

			sess, _ := session.Get("credential", c)
			sess.Options = &sessions.Options{
				Path:     "/",
				MaxAge:   86400 * 7,
				HttpOnly: true,
			}

			sess.Values["user_id"] = user.GetID()
			sess.Values["user_name"] = user.GetUserName()
		}

		tx.Commit()

		return pass, nil
	})
}
