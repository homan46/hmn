package cli

import (
	"fmt"
	"log"
	"os"
	"strings"

	"codeberg.org/rchan/hmn/business"
	"codeberg.org/rchan/hmn/config"
	"codeberg.org/rchan/hmn/helper"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"golang.org/x/term"
)

var resetPasswordCmd = &cobra.Command{
	Use:   "reset",
	Short: "reset password for a user",
	Run: func(cmd *cobra.Command, args []string) {

		configPath, err := cmd.Flags().GetString("config")

		if err != nil {
			log.Fatal("fail to parse config path")
		}

		user, err := cmd.Flags().GetString("user")
		if err != nil {
			log.Fatal("no user provided")
		}

		fmt.Printf("the config path is %s\n", configPath)

		fmt.Printf("please enter new password for %s\n", user)

		bytepw, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			log.Fatal("read line error")
		}

		p1 := strings.TrimSpace(string(bytepw))

		fmt.Printf("please enter new password for %s again\n", user)

		bytepw, err = term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			log.Fatal("read line error")
		}

		p2 := strings.TrimSpace(string(bytepw))

		if p1 == p2 {
			//actual work
			conf := config.LoadConfig(configPath)
			db := helper.OpenDB(conf)
			business := business.NewBusunessLayer(db)

			myContext, tx, err := business.GetContextForSystem()
			if err != nil {
				log.Fatalln(err)
			}
			err = business.User().SetUserPassword(myContext, user, p2)
			if err != nil {
				log.Fatalln(err)
				tx.Rollback()
			}
			tx.Commit()

			fmt.Println("reset password success")
		} else {
			fmt.Println("password mismatch, please try again")
		}

	},
}

func init() {
	set := pflag.NewFlagSet("testset", pflag.ExitOnError)
	set.StringP("user", "u", "", "user to set password")
	resetPasswordCmd.Flags().AddFlagSet(set)

	resetPasswordCmd.Flags().Lookup("user")
	resetPasswordCmd.MarkFlagRequired("user")

}
