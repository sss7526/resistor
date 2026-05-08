package resistor

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"
)

///////////////////////////////////////////////////////////////////////////////
// SMD Encoding Mode
///////////////////////////////////////////////////////////////////////////////

/*
SMDEncodingMode controls how EncodeSMD formats output.

SMDAuto:
    Automatically choose 3-digit or 4-digit encoding.
    If not representable, return error.

SMDStandard:
    Force 3/4-digit encoding only.

SMDEIA96:
    Force EIA‑96 encoding (1% precision format).
*/
type SMDEncodingMode int

const (
	SMDAuto SMDEncodingMode = iota
	SMDStandard
	SMDEIA96
)

///////////////////////////////////////////////////////////////////////////////
// Public Decode API
///////////////////////////////////////////////////////////////////////////////

/*
DecodeSMD decodes a surface-mount resistor marking.

Supported formats:

3-digit:
    "472" → 47 × 10² = 4700Ω

4-digit:
    "4701" → 470 × 10¹ = 4700Ω

R-notation:
    "4R7" → 4.7Ω
    "R47" → 0.47Ω

EIA‑96:
    "01C" → lookup(01) × multiplier(C)

Tolerance is not encoded in SMD markings.
TolerancePct will be left as zero.
*/
func DecodeSMD(marking string) (ResistorSpec, error) {

	var spec ResistorSpec

	m := strings.TrimSpace(strings.ToUpper(marking))
	if m == "" {
		return spec, fmt.Errorf("empty marking")
	}

	// R-notation
	if strings.ContainsRune(m, 'R') {
		val, err := decodeRNotation(m)
		if err != nil {
			return spec, err
		}
		spec.ResistanceOhms = val
		return spec, nil
	}

	// Pure numeric (3 or 4 digit)
	if isAllDigits(m) {
		val, err := decodeNumericSMD(m)
		if err != nil {
			return spec, err
		}
		spec.ResistanceOhms = val
		return spec, nil
	}

	// EIA-96 format (two digits + letter)
	if len(m) == 3 && unicode.IsDigit(rune(m[0])) && unicode.IsDigit(rune(m[1])) && unicode.IsLetter(rune(m[2])) {
		val, err := decodeEIA96(m)
		if err != nil {
			return spec, err
		}
		spec.ResistanceOhms = val
		return spec, nil
	}

	return spec, fmt.Errorf("unsupported SMD marking format")
}

///////////////////////////////////////////////////////////////////////////////
// Public Encode API
///////////////////////////////////////////////////////////////////////////////

/*
EncodeSMD encodes a resistance value into an SMD marking.

Mode behavior:

SMDAuto:
    Attempt standard 3/4 digit encoding first.
    If not possible, error.

SMDStandard:
    Only use 3/4 digit encoding.

SMDEIA96:
    Attempt EIA‑96 encoding.
    If value cannot be represented exactly in E96 series, error.
*/
func EncodeSMD(resistance float64, mode SMDEncodingMode) (string, error) {

	if resistance <= 0 {
		return "", fmt.Errorf("resistance must be positive")
	}

	switch mode {

	case SMDAuto, SMDStandard:
		return encodeStandardSMD(resistance)
	
	case SMDEIA96:
		return encodeEIA96(resistance)

	default:
		return "", fmt.Errorf("unsupported encoding mode")
	}
}

///////////////////////////////////////////////////////////////////////////////
// 3/4 Digit Numeric Decoding
///////////////////////////////////////////////////////////////////////////////

func decodeNumericSMD(m string) (float64, error) {

	if len(m) == 3 {
		// XY Z → (XY) × 10^Z
		base, _ := strconv.Atoi(m[:2])
		exp, _ := strconv.Atoi(string(m[2]))
		return float64(base) * math.Pow(10, float64(exp)), nil
	}

	if len(m) == 4 {
		// XYZ W → (XYZ) × 10^W
		base, _ := strconv.Atoi(m[:3])
		exp, _ := strconv.Atoi(string(m[3]))
		return float64(base) * math.Pow(10, float64(exp)), nil
	}

	return 0, fmt.Errorf("invalid numeric SMD length")
}

///////////////////////////////////////////////////////////////////////////////
// R-Notation Decoding
///////////////////////////////////////////////////////////////////////////////

func decodeRNotation(m string) (float64, error) {

	parts := strings.Split(m, "R")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid R notation")
	}

	left := parts[0]
	right := parts[1]

	if left == "" {
		left = "0"
	}
	if right == "" {
		right = "0"
	}
	
	combined := left + "." + right
	return strconv.ParseFloat(combined, 64)
}

///////////////////////////////////////////////////////////////////////////////
// EIA‑96 Support
///////////////////////////////////////////////////////////////////////////////

var eia96Base = generateESeries(96)

var eia96Multipliers = map[rune]float64{
	'Z': 0.001,
	'Y': 0.01,
	'R': 0.1,
	'X': 1,
	'S': 10,
	'A': 100,
	'B': 1_000,
	'H': 10_000,
	'C': 100_000,
	'D': 1_000_000,
	'E': 10_000_000,
	'F': 100_000_000,
}

func decodeEIA96(m string) (float64, error) {

	code, err := strconv.Atoi(m[:2])
	if err != nil || code < 1 || code > 96 {
		return 0, fmt.Errorf("invalid EIA-96 code")
	}

	base := eia96Base[code-1]

	mult, ok := eia96Multipliers[rune(m[2])]
	if !ok {
		return 0, fmt.Errorf("invalid EIA-96 multiplier")
	}

	return roundToSignificant(base * mult, 6), nil
}

func encodeEIA96(resistance float64) (string, error) {

	// Normalize to decade
	exponent := math.Floor(math.Log10(resistance))
	normalized := resistance / math.Pow(10, exponent)

	// Find closest E96 base
	index := -1
	for i, v := range eia96Base {
		if math.Abs(v-normalized) < 1e-6 {
			index = i
			break
		}
	}

	if index == -1 {
		return "", fmt.Errorf("value not representable in EIA-96 series")
	}

	multChar, ok := findEIA96Multiplier(math.Pow(10, exponent))
	if !ok {
		return "", fmt.Errorf("no valid EIA-96 multiplier")
	}

	return fmt.Sprintf("%02d%c", index + 1, multChar), nil
}

func findEIA96Multiplier(mult float64) (rune, bool) {
	for k, v := range eia96Multipliers {
		if math.Abs(v - mult) < 1e9 {
			return k, true
		}
	}
	return 0, false
}

///////////////////////////////////////////////////////////////////////////////
// Standard Encoding
///////////////////////////////////////////////////////////////////////////////

func encodeStandardSMD(resistance float64) (string, error) {

	for exp := -2; exp <= 9; exp++ {

		scaled := resistance / math.Pow(10, float64(exp))

		if scaled >= 10 && scaled < 100 && math.Mod(scaled, 1) == 0 {
			return fmt.Sprintf("%02d%d", int(scaled), exp), nil
		}

		if scaled >= 100 && scaled < 1000 && math.Mod(scaled, 1) == 0 {
			return fmt.Sprintf("%03d%d", int(scaled), exp), nil
		}
	}

	return "", fmt.Errorf("value not representable in 3/4-digit SMD format")
}

///////////////////////////////////////////////////////////////////////////////
// Helpers
///////////////////////////////////////////////////////////////////////////////

func isAllDigits(s string) bool {
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}