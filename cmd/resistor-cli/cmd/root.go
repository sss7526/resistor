package cmd

import (
	"github.com/spf13/cobra"
)

var jsonOutput bool

var rootCmd = &cobra.Command{
	Use: "resistor-cli",
	Short: "Engineering toolkit for resistor analysis",
	Long: "Resistor CLI provides deterministic standards support, inference, and engineering analysis",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "Output result as JSON")
}