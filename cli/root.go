package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	rootCmd *cobra.Command
)

func init() {

	rootCmd = &cobra.Command{
		Use: "hmn",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("root command")
		},
	}

	rootCmd.PersistentFlags().StringP("config", "c", "config.json", "set config file")

	rootCmd.AddCommand(startServerCmd)
	rootCmd.AddCommand(resetPasswordCmd)

}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
