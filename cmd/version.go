package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

const major = "0"
const minor = "1"
const fix = "0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Describes version.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(fmt.Sprintf("Version: %s.%s.%s", major, minor, fix))
	},
}
