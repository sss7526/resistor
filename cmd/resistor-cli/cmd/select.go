package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/sss7526/resistor"
	"github.com/sss7526/resistor/internal/cli"
)

var (
	selectSeries    string
	selectTolerance float64
	selectRound     string
)

var selectCmd = &cobra.Command{
	Use:   "select [resistance]",
	Short: "Select nearest standard resistor value",
	Args:  cobra.ExactArgs(1),
	Example: `
  # Select nearest standard value (default E24)
  resistor-cli select 487

  # Use specific E-series
  resistor-cli select 487 --series E12

  # Specify tolerance
  resistor-cli select 487 --tol 1

  # Control rounding
  resistor-cli select 487 --round up
  resistor-cli select 487 --round down

  # JSON output
  resistor-cli select 487 --json
`,
	RunE: func(cmd *cobra.Command, args []string) error {

		value, err := parseFloatArg(args[0])
		if err != nil {
			return cli.Respond(jsonOutput, nil, err)
		}
		if value <= 0 {
			return cli.Respond(jsonOutput, nil, fmt.Errorf("resistance must be positive"))
		}

		series, err := parseESeries(selectSeries)
		if err != nil {
			return cli.Respond(jsonOutput, nil, err)
		}

		rounding, err := parseRounding(selectRound)
		if err != nil {
			return cli.Respond(jsonOutput, nil, err)
		}

		req := resistor.SelectionRequest{
			Resistance:   value,
			Series:       series,
			TolerancePct: selectTolerance,
			Rounding:     rounding,
		}

		result, err := resistor.SelectStandardResistor(req)
		if err != nil {
			return cli.Respond(jsonOutput, nil, err)
		}

		if jsonOutput {
			return cli.Respond(jsonOutput, result, nil)
		}

		fmt.Printf("Requested: %-10.6gΩ\n", result.RequestedResistance)
		fmt.Printf("Selected:  %-10.6gΩ\n", result.SelectedResistance)
		cli.PrintBands(result.Bands)

		for _, a := range result.Assumptions {
			fmt.Printf("Note: %s\n", a)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(selectCmd)

	selectCmd.Flags().Float64Var(&selectTolerance, "tol", 0, "Tolerance percentage")
	selectCmd.Flags().StringVar(&selectSeries, "series", "", "Preferred E-series (E3, E6, E12, E24, E48, E96, E192)")
	selectCmd.Flags().StringVar(&selectRound, "round", "", "Rounding mode: nearest, up, down")
}
