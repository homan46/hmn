package controller

import (
	"net/http"

	"codeberg.org/rchan/hmn/business"
	"github.com/labstack/echo/v4"
)

type ViewController struct {
	b business.BusinessLayer
}

func NewViewController(b business.BusinessLayer) *NoteController {
	return &NoteController{
		b: b,
	}
}

func (n *NoteController) GetMainPage(c echo.Context) error {
	csrfToken := c.Get("csrf").(string)

	return c.Render(http.StatusOK, "index.html", map[string]interface{}{
		"csrf_token": csrfToken,
	})
}
