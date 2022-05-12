package midd

import (
	"crypto/rand"
	"log"
	"strings"

	"codeberg.org/rchan/hmn/config"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func NewSess(conf *config.Config) echo.MiddlewareFunc {

	key1 := make([]byte, 32)

	if strings.TrimSpace(conf.Server.CookieSceret) != "" {
		key1 = []byte(strings.TrimSpace(conf.Server.CookieSceret))
	} else {
		_, err := rand.Read(key1)
		if err != nil {
			log.Fatalln(err)
		}
	}

	defaultConfig := session.DefaultConfig
	store := sessions.NewFilesystemStore("")

	store.Codecs = append(store.Codecs, securecookie.CodecsFromPairs(key1)...)

	defaultConfig.Store = store

	return session.MiddlewareWithConfig(defaultConfig)

}
