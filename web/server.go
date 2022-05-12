package web

import (
	"net/http"

	"codeberg.org/rchan/hmn/business"
	"codeberg.org/rchan/hmn/config"
	"codeberg.org/rchan/hmn/web/controller"
	"codeberg.org/rchan/hmn/web/midd"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type a struct {
	name string
}

func New(bl business.BusinessLayer, conf *config.Config) *echo.Echo {
	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: conf.Server.AllowOrigins,
	}))

	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup:    "header:X-XSRF-TOKEN",
		CookieSameSite: http.SameSiteStrictMode,
		CookieHTTPOnly: true,
		CookiePath:     "/",
		Skipper: func(c echo.Context) bool {
			return c.Path() == "/login" && c.Request().Method == http.MethodPost
		},
	}))

	e.Use(midd.NewSess(conf))
	//e.Use(midd.NewAuth(bl))
	//e.Use(midd.NewFakeAuth())
	e.Use(midd.NewDefaultAuth(bl))

	viewRenderer := midd.NewRenderer()
	e.Renderer = viewRenderer

	//e.Static("/", "./public/html")
	e.Static("/js", "./public/js")
	e.Static("/css", "./public/css")
	e.Static("/lib", "./public/lib")

	notec := controller.NewNoteController(bl)
	viewc := controller.NewViewController(bl)
	authc := controller.NewAuthController(bl)

	e.GET("/", viewc.GetMainPage)
	e.GET("/login", viewc.GetLoginPage)
	e.POST("/login", authc.Login)

	v1 := e.Group("/api/v1")

	noteRoute := v1.Group("/note")
	noteRoute.GET("/:id", notec.GetNoteEndpoint)
	noteRoute.GET("", notec.GetAllNoteEndpoint)
	noteRoute.POST("", notec.AddNoteEndpoint)
	noteRoute.PUT("/:id", notec.UpdateNoteEndpoint)
	noteRoute.PATCH("/:id", notec.PatchNoteEndpoint)
	noteRoute.DELETE("/:id", notec.DeleteNoteEndpoint)

	return e
}
