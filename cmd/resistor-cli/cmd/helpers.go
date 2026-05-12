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
	// switch strings.ToUpper(strings.TrimSpace(input)) {
	// case "":
	// 	return 0, nil
	// case "E3":
	// 	return resistor.E3, nil
	// case "E6":
	// 	return resistor.E6, nil
	// case "E12":
	// 	return resistor.E12, nil
	// case "E24":
	// 	return resistor.E24, nil
	// case "E48":
	// 	return resistor.E48, nil
	// case "E96":
	// 	return resistor.E96, nil
	// case "E192":
	// 	return resistor.E192, nil
	// default:
	// 	return 0, fmt.Errorf("invalid E-series: %s", input)
	// }
	return resistor.ParseESeries(input)
}

func parseRounding(input string) (resistor.RoundingMode, error) {
	// switch strings.ToLower(strings.TrimSpace(input)) {
	// case "", "nearest":
	// 	return resistor.RoundNearest, nil
	// case "up":
	// 	return resistor.RoundUp, nil
	// case "down":
	// 	return resistor.RoundDown, nil
	// default:
	// 	return 0, fmt.Errorf("invalid rounding mode: %s", input)
	// }
	return resistor.ParseRoundingMode(input)
}

func parsePackage(input string) (resistor.PackageType, error) {
	// switch strings.TrimSpace(input) {
	// case "":
	// 	return "", nil
	// case "0402":
	// 	return resistor.SMD0402, nil
	// case "0603":
	// 	return resistor.SMD0603, nil
	// case "0805":
	// 	return resistor.SMD0805, nil
	// case "1206":
	// 	return resistor.SMD1206, nil
	// default:
	// 	return "", fmt.Errorf("invalid package: %s", input)
	// }
	return resistor.ParsePackageType(input)
}
