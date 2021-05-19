package cmd

import (
	"fmt"
	"github.com/laupski/online-judge/api"
	"github.com/spf13/cobra"
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Starts the online-judge API server.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {
			if args[0] == "start" {
				api.StartAPI()
				return
			}
		}
		fmt.Println("To start the API server, type start")
	},
}
