package cli

import (
	"fmt"

	"codeberg.org/rchan/hmn/business"
	"codeberg.org/rchan/hmn/config"
	"codeberg.org/rchan/hmn/helper"

	"codeberg.org/rchan/hmn/web"
	"github.com/spf13/cobra"
)

var startServerCmd = &cobra.Command{
	Use:   "start",
	Short: "start listening web server",
	Run: func(cmd *cobra.Command, args []string) {

		configPath, err := cmd.Flags().GetString("config")

		if err != nil {

		}

		fmt.Printf("the config path is %s", configPath)

		conf := config.LoadConfig(configPath)
		db := helper.OpenDB(conf)
		business := business.NewBusunessLayer(db)
		server := web.New(business, conf)

		if conf.Server.UseHttps {
			server.Logger.Fatal(server.StartTLS(conf.Server.ListenOn, conf.Server.TlsCert, conf.Server.TlsKey))
		} else {
			server.Logger.Fatal(server.Start(conf.Server.ListenOn))
		}

		fmt.Println("start")
	},
}
