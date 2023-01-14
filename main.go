package main

import (
	"codeberg.org/rchan/hmn/log"

	"codeberg.org/rchan/hmn/cli"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	log.ZLog.Info("Application Started")
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
