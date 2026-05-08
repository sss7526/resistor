package main

import (
    "os"

    "github.com/sss7526/resistor/cmd/resistor-cli/cmd"
)

func main() {
    if err := cmd.Execute(); err != nil {
        os.Exit(1)
    }
}