package controller

import (
	"net/http"

	"codeberg.org/rchan/hmn/business"
	"github.com/labstack/echo/v4"
)

type ViewController struct {
	b business.BusinessLayer
}

func NewViewController(b business.BusinessLayer) *ViewController {
	return &ViewController{
		b: b,
	}
}

func (n *ViewController) GetMainPage(c echo.Context) error {
	csrfToken := c.Get("csrf").(string)

	return c.Render(http.StatusOK, "index.html", map[string]interface{}{
		"csrf_token": csrfToken,
	})
}
