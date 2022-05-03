package midd

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func NewSess() echo.MiddlewareFunc {
	return session.Middleware(sessions.NewFilesystemStore(""))
}
