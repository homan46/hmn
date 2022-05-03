package cli

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"codeberg.org/rchan/hmn/business"
	"codeberg.org/rchan/hmn/config"
	"codeberg.org/rchan/hmn/data"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"golang.org/x/term"
)

var resetPasswordCmd = &cobra.Command{
	Use:   "reset",
	Short: "reset password for a user",
	Run: func(cmd *cobra.Command, args []string) {

		//configFlag := cmd.Flags().Lookup("config")
		// configFlag.Value.(string)
		// fmt.Println(configFlag.Value)

		configPath, err := cmd.Flags().GetString("config")

		if err != nil {
			log.Fatal("fail to parse config path")
		}

		user, err := cmd.Flags().GetString("user")
		if err != nil {
			log.Fatal("no user provided")
		}

		fmt.Printf("the config path is %s\n", configPath)
		fmt.Printf("the user is %s\n", user)

		writer := bufio.NewWriter(os.Stdout)

		writer.WriteString("please enter new password for " + user + "\n")
		writer.Flush()

		term.ReadPassword(int(os.Stdin.Fd()))

		bytepw, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			log.Fatal("read line error")
		}

		p1 := strings.TrimSpace(string(bytepw))

		writer.WriteString("please enter new password for " + user + " again\n")
		writer.Flush()

		bytepw, err = term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			log.Fatal("read line error")
		}

		p2 := strings.TrimSpace(string(bytepw))

		fmt.Println(p1)
		fmt.Println(p2)

		if p1 == p2 {
			//actual work
			conf := config.LoadConfig(configPath)
			db := data.OpenDB(conf)
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
			writer.WriteString("reset password success")
			writer.Flush()
		} else {
			writer.WriteString("password mismatch, no change \n")
			writer.Flush()
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

func readPassword() {

}
