package cmd

import (
    "fmt"
    "strings"

    "github.com/spf13/cobra"

    "github.com/sss7526/resistor"
    "github.com/sss7526/resistor/internal/cli"
)

var (
    inferBands string
    inferSMD   string
    inferLen   float64
    inferBody  string
)

var inferCmd = &cobra.Command{
	Use: "infer",
	Short: "Infer resistor properties from observed data",
	RunE: func(cmd *cobra.Command, args []string) error {

		obs := resistor.ObservedResistor{
			LengthMM: inferLen,
			BodyColor: resistor.Color(strings.ToLower(inferBody)),
		}

		if inferBands != "" {
			bands, err := parseBands(inferBands)
			if err != nil {
				return cli.Respond(jsonOutput, nil, err)
			}
			obs.Bands = bands
		}

		if inferSMD != "" {
			obs.Marking = inferSMD
		}

		result, err := resistor.InferResistor(obs)
		if err != nil {
			return cli.Respond(jsonOutput, nil, err)
		}

		if jsonOutput {
			return cli.Respond(jsonOutput, result, nil)
		}

		fmt.Printf("Resistance: %.6gΩ\n", result.Spec.ResistanceOhms)
        fmt.Printf("Tolerance:  %.2f%%\n", result.Spec.TolerancePct)
        fmt.Printf("Power:      %.3gW\n", result.Spec.PowerWatts)
        fmt.Printf("Confidence: %.2f\n", result.Meta.Confidence)

        for _, a := range result.Meta.Assumptions {
            fmt.Printf("Note: %s\n", a)
        }

        return nil
	},
}

func init() {
	rootCmd.AddCommand(inferCmd)

	inferCmd.Flags().StringVar(&inferBands, "bands", "", "Comma-separated band colors")
	inferCmd.Flags().StringVar(&inferSMD, "smd", "", "SMD marking")
	inferCmd.Flags().Float64Var(&inferLen, "length", 0, "Length in mm")
	inferCmd.Flags().StringVar(&inferBody, "body", "", "Body color")
}