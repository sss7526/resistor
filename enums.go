package resistor

import (
	"fmt"
	"strings"
)

///////////////////////////////////////////////////////////////////////////////
// ESeries
///////////////////////////////////////////////////////////////////////////////

// String returns the canonical string representation of an E-series.
func (e ESeries) String() string {
	switch e {
	case E3:
		return "E3"
	case E6:
		return "E6"
	case E12:
		return "E12"
	case E24:
		return "E24"
	case E48:
		return "E48"
	case E96:
		return "E96"
	case E192:
		return "E192"
	default:
		return "Unknown"
	}
}

// ParseESeries converts a string into an ESeries value.
func ParseESeries(input string) (ESeries, error) {
	switch strings.ToUpper(strings.TrimSpace(input)) {
	case "E3":
		return E3, nil
	case "E6":
		return E6, nil
	case "E12":
		return E12, nil
	case "E24":
		return E24, nil
	case "E48":
		return E48, nil
	case "E96":
		return E96, nil
	case "E192":
		return E192, nil
	default:
		return 0, fmt.Errorf("invalid E-series: %s", input)
	}
}

// AllESeries returns the supported E-series values.
func AllESeries() []ESeries {
	return []ESeries{E3, E6, E12, E24, E48, E96, E192}
}

///////////////////////////////////////////////////////////////////////////////
// RoundingMode
///////////////////////////////////////////////////////////////////////////////

// String returns the canonical string representation of a rounding mode.
func (r RoundingMode) String() string {
	switch r {
	case RoundNearest:
		return "Nearest"
	case RoundUp:
		return "Up"
	case RoundDown:
		return "Down"
	default:
		return "Unknown"
	}
}

// ParseRoundingMode converts a string into a RoundingMode.
func ParseRoundingMode(input string) (RoundingMode, error) {
	switch strings.ToLower(strings.TrimSpace(input)) {
	case "", "nearest":
		return RoundNearest, nil
	case "up":
		return RoundUp, nil
	case "down":
		return RoundDown, nil
	default:
		return 0, fmt.Errorf("invalid rounding mode: %s", input)
	}
}

// AllRoundingModes returns all supported rounding modes.
func AllRoundingModes() []RoundingMode {
	return []RoundingMode{RoundNearest, RoundUp, RoundDown}
}

///////////////////////////////////////////////////////////////////////////////
// PackageType
///////////////////////////////////////////////////////////////////////////////

// String returns the canonical string representation of a package type.
func (p PackageType) String() string {
	switch p {
	case ThroughHole:
		return "through_hole"
	case SMD:
		return "smd"
	case SMD0402:
		return "0402"
	case SMD0603:
		return "0603"
	case SMD0805:
		return "0805"
	case SMD1206:
		return "1206"
	case Axial:
		return "axial"
	case Radial:
		return "radial"
	default:
		return "unknown"
	}
}

// ParsePackageType converts a string into a PackageType.
func ParsePackageType(input string) (PackageType, error) {
	switch strings.ToLower(strings.TrimSpace(input)) {
	case "":
		return "", nil
	case "through_hole":
		return ThroughHole, nil
	case "smd":
		return SMD, nil
	case "0402":
		return SMD0402, nil
	case "0603":
		return SMD0603, nil
	case "0805":
		return SMD0805, nil
	case "1206":
		return SMD1206, nil
	case "axial":
		return Axial, nil
	case "radial":
		return Radial, nil
	default:
		return "", fmt.Errorf("invalid package type: %s", input)
	}
}

// AllPackageTypes returns the commonly selectable package types.
// UnknownPKG is excluded; it is used internally when a type cannot be determined.
func AllPackageTypes() []PackageType {
	return []PackageType{
		ThroughHole,
		SMD0402,
		SMD0603,
		SMD0805,
		SMD1206,
		Axial,
		Radial,
	}
}

///////////////////////////////////////////////////////////////////////////////
// Color Enumerations
///////////////////////////////////////////////////////////////////////////////

// The following functions expose canonical color lists for UI layers.
// These are intentionally ordered for deterministic presentation.
// UI layers (CLI, TUI, WASM) should rely on these functions rather
// than hardcoding color lists.

// DigitColors returns colors valid for significant digit bands (0–9).
func DigitColors() []Color {
	return []Color{
		Black,
		Brown,
		Red,
		Orange,
		Yellow,
		Green,
		Blue,
		Violet,
		Grey,
		White,
	}
}

// MultiplierColors returns valid multiplier band colors.
func MultiplierColors() []Color {
	return []Color{
		Black,
		Brown,
		Red,
		Orange,
		Yellow,
		Green,
		Blue,
		Violet,
		Grey,
		White,
		Gold,
		Silver,
	}
}

// ToleranceColors returns valid tolerance band colors.
func ToleranceColors() []Color {
	return []Color{
		Brown,
		Red,
		Green,
		Blue,
		Violet,
		Grey,
		Gold,
		Silver,
		None,
	}
}

// TempCoeffColors returns valid 6th band temperature coefficient colors.
func TempCoeffColors() []Color {
	return []Color{
		Brown,
		Red,
		Orange,
		Yellow,
		Blue,
		Violet,
	}
}
