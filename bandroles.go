package resistor

import "fmt"

// This file models structural band roles defined by IEC 60062.
//
// It provides:
//
//   - Ordered band role definitions
//   - Valid color enumeration per role
//   - Canonical structural rules for 4/5/6-band resistors
//
// UI layers (CLI, TUI, WASM) must rely on this API
// instead of hardcoding band structure logic.

///////////////////////////////////////////////////////////////////////////////
// Band Roles (IEC 60062 Structural Modeling)
///////////////////////////////////////////////////////////////////////////////

/*
BandRole describes the semantic role of a resistor band position.

Roles are derived from IEC 60062 band structure definitions.

These roles are used to:

  - Determine valid colors for a band position
  - Drive UI selection logic
  - Prevent duplicated band structure logic in frontends
*/
type BandRole int

const (
	RoleDigit BandRole = iota
	RoleMultiplier
	RoleTolerance
	RoleTempCoeff
)

func (r BandRole) String() string {
	switch r {
	case RoleDigit:
		return "Digit"
	case RoleMultiplier:
		return "Multiplier"
	case RoleTolerance:
		return "Tolerance"
	case RoleTempCoeff:
		return "TempCoeff"
	default:
		return "Unknown"
	}
}

/*
BandRolesForCount returns the ordered band roles for a resistor
with the specified band count.

Supported band counts:

  - 4 bands
  - 5 bands
  - 6 bands

IEC 60062 definitions:

4-band:

	Digit, Digit, Multiplier, Tolerance

5-band:

	Digit, Digit, Digit, Multiplier, Tolerance

6-band:

	Digit, Digit, Digit, Multiplier, Tolerance, TempCoeff
*/
func BandRolesForCount(count int) ([]BandRole, error) {

	switch count {

	case 4:
		return []BandRole{
			RoleDigit,
			RoleDigit,
			RoleMultiplier,
			RoleTolerance,
		}, nil

	case 5:
		return []BandRole{
			RoleDigit,
			RoleDigit,
			RoleDigit,
			RoleMultiplier,
			RoleTolerance,
		}, nil

	case 6:
		return []BandRole{
			RoleDigit,
			RoleDigit,
			RoleDigit,
			RoleMultiplier,
			RoleTolerance,
			RoleTempCoeff,
		}, nil

	default:
		return nil, fmt.Errorf("unsupported band count: %d", count)
	}
}

/*
ValidColorsForRole returns the valid colors for a given band role.

This function exposes canonical color sets for UI layers
and ensures no IEC duplication outside the core.
*/
func ValidColorsForRole(role BandRole) []Color {

	switch role {

	case RoleDigit:
		return DigitColors()

	case RoleMultiplier:
		return MultiplierColors()

	case RoleTolerance:
		return ToleranceColors()

	case RoleTempCoeff:
		return TempCoeffColors()

	default:
		return nil
	}
}
