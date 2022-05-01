package main

import (
	"log"

	"codeberg.org/rchan/hmn/business"
	"codeberg.org/rchan/hmn/config"
	"codeberg.org/rchan/hmn/data"
	"codeberg.org/rchan/hmn/web"
	_ "github.com/mattn/go-sqlite3"
)

func init() {
	log.SetFlags(log.Llongfile | log.LUTC)
}

func main() {
	conf := config.LoadConfig("./config.json")
	db := data.OpenDB(conf)
	business := business.NewBusunessLayer(db)
	server := web.New(business, conf)

	if conf.Server.UseHttps {
		server.Logger.Fatal(server.StartTLS(conf.Server.ListenOn, conf.Server.TlsCert, conf.Server.TlsKey))
	} else {
		server.Logger.Fatal(server.Start(conf.Server.ListenOn))
	}

}
