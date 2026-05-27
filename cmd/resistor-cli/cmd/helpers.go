package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/sss7526/resistor"
)

func parseFloatArg(arg string) (float64, error) {
	v, err := strconv.ParseFloat(arg, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid numeric value: %s", arg)
	}
	return v, nil
}

func parseBands(input string) ([]resistor.Color, error) {
	parts := strings.Split(input, ",")
	var bands []resistor.Color

	for _, p := range parts {
		color := resistor.Color(strings.ToLower(strings.TrimSpace(p)))

		_, digitOK := resistor.DigitValue[color]
		_, multOK := resistor.MultiplierValue[color]
		_, tolOK := resistor.ToleranceValue[color]
		_, tempOK := resistor.TempCoeffValue[color]

		if !(digitOK || multOK || tolOK || tempOK) {
			return nil, fmt.Errorf("invalid band color: %s", p)
		}

		bands = append(bands, color)
	}

	return bands, nil
}

func parseESeries(input string) (resistor.ESeries, error) {
	if strings.TrimSpace(input) == "" {
		return 0, nil
	}
	return resistor.ParseESeries(input)
}

func parseRounding(input string) (resistor.RoundingMode, error) {
	return resistor.ParseRoundingMode(input)
}

func parsePackage(input string) (resistor.PackageType, error) {
	return resistor.ParsePackageType(input)
}
