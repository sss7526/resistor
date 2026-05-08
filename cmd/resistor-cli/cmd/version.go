package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var version = "v0.1.0"

var versionCmd = &cobra.Command{
	Use: "version",
	Short: "Print CLI version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}