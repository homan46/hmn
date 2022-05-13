package midd

import (
	"net/http"

	"codeberg.org/rchan/hmn/business"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewAuth(bl business.BusinessLayer) echo.MiddlewareFunc {

	return middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {

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

func NewFakeAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sess, _ := session.Get("credential", c)
			sess.Options = &sessions.Options{
				Path:     "/",
				MaxAge:   86400 * 7,
				HttpOnly: true,
			}

			sess.Values["user_id"] = 2
			sess.Values["user_name"] = "admin"
			return next(c)
		}
	}
}

func NewDefaultAuth(bl business.BusinessLayer) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {

		return func(c echo.Context) error {

			sess, err := session.Get("credential", c)
			if err != nil {
				return err
			}

			//first time
			if sess.IsNew {
				err := sess.Save(c.Request(), c.Response())
				if err != nil {
					return err
				}
			}

			_, exists := sess.Values["user_id"]

			if exists {
				//redirect to main page if trying to go to login page after login
				if c.Path() == "/login" {
					return c.Redirect(http.StatusFound, "/")
				}
				return next(c)
			} else {
				//always allow access to login page
				if c.Path() == "/login" || c.Path() == "/api/v1/session" {
					return next(c)
				}
				//no user_id means user is not login
				return c.Redirect(http.StatusFound, "/login")

			}

		}
	}
}
