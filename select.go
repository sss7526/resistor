package resistor

import (
	"fmt"
)

///////////////////////////////////////////////////////////////////////////////
// Selection Models
///////////////////////////////////////////////////////////////////////////////

/*
SelectionRequest represents an engineering intent to select
a standard resistor.

Zero-value fields are treated as unspecified and will
be replaced with sensible defaults.
*/
type SelectionRequest struct {
	// Desired resistance in ohms.
	Resistance float64

	// Optional preferred E-series.
	// If zero-value, defaults to E24.
	Series ESeries

	// Optional tolerance percentage.
	// If zero-value, defaults to 5%.
	TolerancePct float64

	// Optional rounding mode.
	// If zero-value, defaults to RoundNearest.
	Rounding RoundingMode
}

/*
SelectionResult represents the fully resolved resistor selection.

This includes:
  - The originally requested resistance
  - The snapped standard resistance
  - The E-series used
  - The tolerance applied
  - The rounding mode used
  - The resulting color bands
  - Explicit assumptions made during selection
*/
type SelectionResult struct {
	RequestedResistance float64
	SelectedResistance  float64

	Series       ESeries
	TolerancePct float64
	Rounding     RoundingMode

	Bands []Color

	Assumptions []string
}

///////////////////////////////////////////////////////////////////////////////
// Public API
///////////////////////////////////////////////////////////////////////////////

/*
SelectStandardResistor performs a full forward transformation:

	Engineering intent → Standard value → Color bands

Process:

1. Validate requested resistance.
2. Apply defaults where fields are unspecified.
3. Snap resistance to preferred E-series value.
4. Encode selected value into color bands.
5. Return structured result with explicit assumptions.

Defaults:

	Series        → E24
	Tolerance     → 5%
	Rounding Mode → RoundNearest

This function is deterministic and does not perform inference.
*/
func SelectStandardResistor(req SelectionRequest) (SelectionResult, error) {

	var result SelectionResult
	var assumptions []string

	if req.Resistance <= 0 {
		return result, fmt.Errorf("resistance must be positive")
	}

	result.RequestedResistance = req.Resistance

	// ---------------------------------------------------------------------
	// Apply Defaults
	// ---------------------------------------------------------------------

	series := req.Series
	if series == 0 {
		series = E24
		assumptions = append(assumptions, "Series defaulted to E24")
	}

	tolerance := req.TolerancePct
	if tolerance == 0 {
		tolerance = 5.0
		assumptions = append(assumptions, "Tolerance defaulted to ±5%")
	}

	rounding := req.Rounding
	if rounding == RoundingUnspecified {
		rounding = RoundNearest
		assumptions = append(assumptions, "Rounding mode defaulted to RoundNearest")
	}

	result.Series = series
	result.TolerancePct = tolerance
	result.Rounding = rounding

	// ---------------------------------------------------------------------
	// Snap to Standard Value
	// ---------------------------------------------------------------------

	selected, err := NearestStandard(req.Resistance, series, rounding)
	if err != nil {
		return result, err
	}

	if selected != req.Resistance {
		assumptions = append(
			assumptions,
			fmt.Sprintf(
				"Resistance snapped from %.6gΩ to %.6gΩ",
				req.Resistance,
				selected,
			),
		)
	}

	result.SelectedResistance = selected

	// ---------------------------------------------------------------------
	// Encode Color Bands
	// ---------------------------------------------------------------------

	bands, err := EncodeBands(ResistorSpec{ResistanceOhms: selected, TolerancePct: tolerance})
	if err != nil {
		return result, err
	}

	result.Bands = bands
	result.Assumptions = assumptions

	return result, nil
}
