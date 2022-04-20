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
	server := web.New(business)
	server.Logger.Debug()
	server.Logger.Fatal(server.Start(":8080"))
}
