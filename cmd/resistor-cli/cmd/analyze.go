package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/sss7526/resistor"
	"github.com/sss7526/resistor/internal/cli"
)

var (
	anR   float64
	anV   float64
	anI   float64
	anP   float64
	anTol float64
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Perform engineering analysis on a resistor",
	Example: `
  # Voltage-driven analysis
  resistor-cli analyze --r 100 --v 10 --pwr 0.5

  # Current-driven analysis
  resistor-cli analyze --r 50 --i 0.2

  # Include tolerance for worst-case bounds
  resistor-cli analyze --r 100 --v 10 --pwr 0.5 --tol 5

  # JSON output
  resistor-cli analyze --r 100 --v 10 --pwr 0.5 --json
`,
	RunE: func(cmd *cobra.Command, args []string) error {

		if anR <= 0 {
			return cli.Respond(jsonOutput, nil, fmt.Errorf("resistance (--r) must be positive"))
		}

		spec := resistor.ResistorSpec{
			ResistanceOhms: anR,
			PowerWatts:     anP,
			TolerancePct:   anTol,
		}

		input := resistor.AnalysisInput{
			Spec:           spec,
			AppliedVoltage: anV,
			AppliedCurrent: anI,
		}

		report, err := resistor.AnalyzeResistor(input)
		if err != nil {
			return cli.Respond(jsonOutput, nil, err)
		}

		if jsonOutput {
			return cli.Respond(jsonOutput, report, nil)
		}

		fmt.Printf("Voltage Drop: %-10.6gV\n", report.VoltageDrop)
		fmt.Printf("Current:      %-10.6gA\n", report.Current)
		fmt.Printf("Power:        %-10.6gW\n", report.PowerDissipation)

		if report.DeratedSafePower > 0 {
			fmt.Printf("Derated Safe: %-10.6gW\n", report.DeratedSafePower)
		}

		if report.WorstCaseResistanceMin > 0 {
			fmt.Printf("WC R Min:     %-10.6gΩ\n", report.WorstCaseResistanceMin)
			fmt.Printf("WC R Max:     %-10.6gΩ\n", report.WorstCaseResistanceMax)
		}

		for _, w := range report.Warnings {
			fmt.Printf("[%s] %s\n", w.Level, w.Message)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(analyzeCmd)

	analyzeCmd.Flags().Float64Var(&anR, "r", 0, "Resistance in ohms")
	analyzeCmd.Flags().Float64Var(&anV, "v", 0, "Applied voltage")
	analyzeCmd.Flags().Float64Var(&anI, "i", 0, "Applied current")
	analyzeCmd.Flags().Float64Var(&anP, "pwr", 0, "Rated power (W)")
	analyzeCmd.Flags().Float64Var(&anTol, "tol", 0, "Tolerance (%)")
}
