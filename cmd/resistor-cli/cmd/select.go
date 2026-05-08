package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/sss7526/resistor"
	"github.com/sss7526/resistor/internal/cli"
)

var (
	selectSeries string
	selectTolerance float64
)

var selectCmd = &cobra.Command{
	Use: "select [resistance]",
	Short: "Select nearest standard resistor value",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		value, err := parseFloatArg(args[0])
		if err != nil {
			return cli.Respond(jsonOutput, nil, err)
		}

		series, err := parseESeries(selectSeries)
		if err != nil {
			return cli.Respond(jsonOutput, nil, err)
		}

		req := resistor.SelectionRequest{
			Resistance:   value,
			Series:       series,
			TolerancePct: selectTolerance,
		}

		result, err := resistor.SelectStandardResistor(req)
		if err != nil {
			return cli.Respond(jsonOutput, nil, err)
		}

		if jsonOutput {
			return cli.Respond(jsonOutput, result, nil)
		}

		fmt.Printf("Requested: %.6gΩ\n", result.RequestedResistance)
		fmt.Printf("Selected:  %.6gΩ\n", result.SelectedResistance)
		fmt.Printf("Bands:     %v\n", result.Bands)

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
}