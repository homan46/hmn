package cli

import (
	"fmt"

	"codeberg.org/rchan/hmn/log"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	rootCmd *cobra.Command
)

func init() {

	rootCmd = &cobra.Command{
		Use: "hmn",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("please use -h for more options")
		},
	}

	rootCmd.PersistentFlags().StringP("config", "c", "config.json", "set config file")

	rootCmd.AddCommand(startServerCmd)
	rootCmd.AddCommand(resetPasswordCmd)

}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.ZLog.Panic("root command end with error", zap.Error(err))
	}
}
