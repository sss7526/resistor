package cmd

import (
    "fmt"

    "github.com/spf13/cobra"

    "github.com/sss7526/resistor"
    "github.com/sss7526/resistor/internal/cli"
)

var smdCmd = &cobra.Command{
    Use:   "smd",
    Short: "SMD encode/decode operations",
}

var smdDecodeCmd = &cobra.Command{
    Use:   "decode [marking]",
    Short: "Decode SMD resistor marking",
    Args:  cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {

        spec, err := resistor.DecodeSMD(args[0])
        if err != nil {
			return cli.Respond(jsonOutput, nil, err)
        }

        if jsonOutput {
			return cli.Respond(jsonOutput, spec, nil)
        }

        fmt.Printf("Resistance: %.6gΩ\n", spec.ResistanceOhms)
        return nil
    },
}

var smdEncodeCmd = &cobra.Command{
    Use:   "encode [resistance]",
    Short: "Encode resistance into SMD marking",
    Args:  cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {

        value, err := parseFloatArg(args[0])
        if err != nil {
            return err
        }

        code, err := resistor.EncodeSMD(value, resistor.SMDAuto)
        if err != nil {
            return cli.Respond(jsonOutput, nil, err)
        }

        if jsonOutput {
			return cli.Respond(jsonOutput, map[string]string{
                "marking": code,
            }, nil)
        }

        fmt.Printf("SMD Marking: %s\n", code)
        return nil
    },
}

func init() {
    rootCmd.AddCommand(smdCmd)
    smdCmd.AddCommand(smdDecodeCmd)
    smdCmd.AddCommand(smdEncodeCmd)
}