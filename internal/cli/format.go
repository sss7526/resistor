package cli

import (
	"fmt"
	"github.com/sss7526/resistor"
	"strings"
)

func PrintHeader(title string) {
	fmt.Printf("=== %s ===\n", title)
}

func PrintBands(bands []resistor.Color) {
	fmt.Println("Bands:")
	for _, b := range bands {
		fmt.Printf("  %s\n", strings.Title(string(b)))
	}
}
