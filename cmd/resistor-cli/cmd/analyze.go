package cmd

import (
    "fmt"

    "github.com/spf13/cobra"

    "github.com/sss7526/resistor"
    "github.com/sss7526/resistor/internal/cli"
)

var (
    anR float64
    anV float64
    anI float64
    anP float64
    anTol float64
)

var analyzeCmd = &cobra.Command{
    Use:   "analyze",
    Short: "Perform engineering analysis on a resistor",
    RunE: func(cmd *cobra.Command, args []string) error {

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

        fmt.Printf("Voltage Drop: %.6gV\n", report.VoltageDrop)
        fmt.Printf("Current:      %.6gA\n", report.Current)
        fmt.Printf("Power:        %.6gW\n", report.PowerDissipation)

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