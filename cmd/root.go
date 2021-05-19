package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "online-judge",
	Short: "Online Judge is an online remote code judge system",
	Long: `A proof of concept of an online remote code judge system 
                written by laupski in Go.
                Complete documentation is available at https://github.com/laupski/online-judge`,
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(apiCmd)
	rootCmd.AddCommand(judgeCmd)
}

// Execute runs the command line input.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
