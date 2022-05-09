package midd

import (
	"crypto/rand"
	"log"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func NewSess() echo.MiddlewareFunc {
	key1 := make([]byte, 32)
	_, err := rand.Read(key1)
	if err != nil {
		log.Fatalln(err)
	}

	defaultConfig := session.DefaultConfig
	store := sessions.NewFilesystemStore("")

	store.Codecs = append(store.Codecs, securecookie.CodecsFromPairs(key1)...)

	defaultConfig.Store = store

	return session.MiddlewareWithConfig(defaultConfig)

}
