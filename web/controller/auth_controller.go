package controller

import (
	"net/http"

	"codeberg.org/rchan/hmn/business"
	"codeberg.org/rchan/hmn/business/service"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type AuthController struct {
	b business.BusinessLayer
}

func NewAuthController(b business.BusinessLayer) *AuthController {
	return &AuthController{
		b: b,
	}
}

func (n *AuthController) Login(c echo.Context) error {
	userName := c.FormValue("userName")
	pass := c.FormValue("password")
	csrfInput := c.FormValue("csrf")

	//TODO: better way for csrf with form?
	cc, err := c.Cookie("_csrf")

	csrfToken := cc.Value
	if csrfInput != csrfToken {
		return echo.ErrBadRequest
	}

	/*
		csrfInput := c.FormValue("csrf")

		csrfToken := c.Get("csrf").(string)

		if csrfInput != csrfToken {
			return echo.ErrBadRequest
		}
	*/
	mycontext, tx, err := n.b.GetContextForSystem()
	if err != nil {
		return err
	}

	ok, err := n.b.User().CheckUserPassword(mycontext, userName, pass)
	if err != nil {
		tx.Rollback()

		if err == service.ErrUserNotExist {
			return c.Redirect(http.StatusFound, "/login")
		}
		return err
	}

	if ok {
		user, err := n.b.User().GetUserByUserName(mycontext, userName)
		if err != nil {
			tx.Rollback()
			return err
		}

		sess, _ := session.Get("credential", c)
		sess.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   86400 * 2,
			HttpOnly: true,
		}

		sess.Values["user_id"] = user.GetID()
		sess.Values["user_name"] = user.GetUserName()

		err = sess.Save(c.Request(), c.Response())
		if err != nil {
			tx.Rollback()
			return err
		}

		tx.Commit()
		return c.Redirect(http.StatusFound, "/")
	}

	tx.Commit()
	return c.Redirect(http.StatusFound, "/login")
}
