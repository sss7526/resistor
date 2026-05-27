package cli

import (
	"fmt"

	"github.com/sss7526/resistor"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var titler = cases.Title(language.Und)

func PrintHeader(title string) {
	fmt.Printf("=== %s ===\n", title)
}

func PrintBands(bands []resistor.Color) {
	fmt.Println("Bands:")
	for _, b := range bands {
		fmt.Printf("  %s\n", titler.String(string(b)))
	}
}
