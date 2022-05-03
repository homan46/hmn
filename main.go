package main

import (
	"log"

	"codeberg.org/rchan/hmn/cli"
	_ "github.com/mattn/go-sqlite3"
)

func init() {
	log.SetFlags(log.Llongfile | log.LUTC)
}

func main() {

	cli.Execute()

	// conf := config.LoadConfig("./config.json")
	// db := data.OpenDB(conf)
	// business := business.NewBusunessLayer(db)
	// server := web.New(business, conf)

	// if conf.Server.UseHttps {
	// 	server.Logger.Fatal(server.StartTLS(conf.Server.ListenOn, conf.Server.TlsCert, conf.Server.TlsKey))
	// } else {
	// 	server.Logger.Fatal(server.Start(conf.Server.ListenOn))
	// }

}
