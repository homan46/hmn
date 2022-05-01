package web

import (
	"codeberg.org/rchan/hmn/business"
	"codeberg.org/rchan/hmn/web/controller"
	"github.com/labstack/echo/v4"
)

type a struct {
	name string
}

func New(bl business.BusinessLayer) *echo.Echo {
	e := echo.New()

	e.Static("/", "./public/html")
	e.Static("/js", "./public/js")
	e.Static("/css", "./public/css")
	e.Static("/lib", "./public/lib")

	notec := controller.NewNoteController(bl)

	noteRoute := v1.Group("/note")
	noteRoute.GET("/:id", notec.GetNoteEndpoint)
	noteRoute.GET("", notec.GetAllNoteEndpoint)
	noteRoute.POST("", notec.AddNoteEndpoint)
	noteRoute.PUT("/:id", notec.UpdateNoteEndpoint)
	noteRoute.PATCH("/:id", notec.PatchNoteEndpoint)
	noteRoute.DELETE("/:id", notec.DeleteNoteEndpoint)

	return e
}
