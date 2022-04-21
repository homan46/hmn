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

	v1 := e.Group("/api/v1")
	v1.GET("/ping", func(c echo.Context) error {
		return c.JSON(200,
			struct {
				Name string `json:"name"`
				Age  int    `json:"age"`
			}{"tom", 10})
	})

	notec := controller.NewNoteController(bl)

	noteRoute := v1.Group("/note")
	noteRoute.GET("/:id", notec.GetNoteEndpoint)
	noteRoute.GET("", notec.GetAllNoteEndpoint)
	noteRoute.POST("", notec.AddNoteEndpoint)
	noteRoute.PUT("", notec.UpdateNoteEndpoint)

	return e
}
