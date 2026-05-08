package cmd

import (
	"github.com/spf13/cobra"
)

var jsonOutput bool

var rootCmd = &cobra.Command{
	Use:   "resistor-cli",
	Short: "Engineering toolkit for resistor analysis",
	Long:  "Resistor CLI provides deterministic standards support, inference, and engineering analysis",
	Example: `
  # Select a standard resistor
  resistor-cli select 487

  # Infer resistor from color bands
  resistor-cli infer --bands brown,black,red,gold

  # Analyze resistor power dissipation
  resistor-cli analyze --r 100 --v 10 --pwr 0.5

  # Decode SMD marking
  resistor-cli smd decode 472
`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "Output result as JSON")
}
