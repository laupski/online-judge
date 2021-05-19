package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var judgeCmd = &cobra.Command{
	Use:   "judge",
	Short: "Starts the online-judge judge server.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {
			if args[0] == "start" {

				return
			}
		}
		fmt.Println("To start the Judge server, type start")
	},
}
