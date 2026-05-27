package resistor

import (
	"fmt"
	"math"
)

///////////////////////////////////////////////////////////////////////////////
// Warning System
///////////////////////////////////////////////////////////////////////////////

// WarningLevel indicates severity of an analysis warning
type WarningLevel string

const (
	WarningInfo    WarningLevel = "info"
	WarningCaution WarningLevel = "caution"
	WarningDanger  WarningLevel = "danger"
)

// AnalysisWarning represents a structured engineering warning.
type AnalysisWarning struct {
	Level   WarningLevel
	Message string
}

///////////////////////////////////////////////////////////////////////////////
// Analysis Models
///////////////////////////////////////////////////////////////////////////////

// AnalysisInput contains electrical conditions for resistor analysis.
//
// Either AppliedVoltage or AppliedCurrent (or both) may be provided.
// If both are provided, consistency is checked.
type AnalysisInput struct {
	Spec ResistorSpec

	AppliedVoltage float64
	AppliedCurrent float64
}

// AnalysisReport contains deterministic electrical analysis results.
//
// DeratedSafePower, WorstCaseResistanceMin, and WorstCaseResistanceMax are
// pointer fields. A nil value means the field was not computed (required input
// was absent). A non-nil pointer to zero is a legitimately computed zero value
// (e.g. WorstCaseResistanceMin at 100% tolerance).
type AnalysisReport struct {
	PowerDissipation float64
	VoltageDrop      float64
	Current          float64

	DeratedSafePower *float64 `json:"DeratedSafePower,omitempty"`

	WorstCaseResistanceMin *float64 `json:"WorstCaseResistanceMin,omitempty"`
	WorstCaseResistanceMax *float64 `json:"WorstCaseResistanceMax,omitempty"`

	Warnings []AnalysisWarning
}

///////////////////////////////////////////////////////////////////////////////
// Public API
///////////////////////////////////////////////////////////////////////////////

// AnalyzeResistor performs deterministic engineering analysis of a resistor
// under specified electrical conditions.
//
// It computes:
//   - Power dissipation
//   - Voltage drop
//   - Current
//   - Derated safe power (50% rule)
//   - Worst-case resistance bounds
//
// The function does not error for missing optional inputs.
// Instead, it produces structured warnings.
func AnalyzeResistor(input AnalysisInput) (AnalysisReport, error) {

	var report AnalysisReport
	var warnings []AnalysisWarning

	spec := input.Spec

	if spec.ResistanceOhms <= 0 {
		return report, fmt.Errorf("resistance must be positive for analysis")
	}

	R := spec.ResistanceOhms
	V := input.AppliedVoltage
	I := input.AppliedCurrent

	// ---------------------------------------------------------
	// Electrical Condition Resolution
	// ---------------------------------------------------------

	if V > 0 && I > 0 {
		// Both provided - check consistency
		expectedV := I * R
		if math.Abs(expectedV-V)/math.Max(math.Abs(V), 1e-12) > 0.01 {
			warnings = append(warnings, AnalysisWarning{
				Level:   WarningCaution,
				Message: "Applied voltage and current inconsistent with Ohm's Law for given resistance",
			})
		}
		report.Current = I
		report.VoltageDrop = V

	} else if V > 0 {
		report.VoltageDrop = V
		report.Current = V / R

	} else if I > 0 {
		report.Current = I
		report.VoltageDrop = I * R

	} else {
		warnings = append(warnings, AnalysisWarning{
			Level:   WarningInfo,
			Message: "No applied voltage or current provided: power dissipation not computed",
		})
	}

	// ---------------------------------------------------------
	// Power Dissipation
	// ---------------------------------------------------------

	if report.Current > 0 {
		report.PowerDissipation = report.Current * report.Current * R
	}

	// ---------------------------------------------------------
	// Derating (50% rule)
	// ---------------------------------------------------------

	if spec.PowerWatts > 0 {
		dsp := 0.5 * spec.PowerWatts
		report.DeratedSafePower = &dsp

		if report.PowerDissipation > spec.PowerWatts {
			warnings = append(warnings, AnalysisWarning{
				Level:   WarningDanger,
				Message: "Power dissipation exceeds rated power",
			})
		} else if report.PowerDissipation > *report.DeratedSafePower {
			warnings = append(warnings, AnalysisWarning{
				Level:   WarningCaution,
				Message: "Power dissipation exceeds recommended 50% derated threshold",
			})
		}

	} else {
		warnings = append(warnings, AnalysisWarning{
			Level:   WarningInfo,
			Message: "Power rating unknown; derating analysis unavailable",
		})
	}

	// ---------------------------------------------------------
	// Tolerance Worst-Case Bounds
	// ---------------------------------------------------------

	if spec.TolerancePct > 0 {
		tol := spec.TolerancePct / 100.0
		wcMin := R * (1 - tol)
		wcMax := R * (1 + tol)
		report.WorstCaseResistanceMin = &wcMin
		report.WorstCaseResistanceMax = &wcMax
	} else {
		warnings = append(warnings, AnalysisWarning{
			Level:   WarningInfo,
			Message: "Tolerance unknown; worst-case resistance not computed",
		})
	}

	report.Warnings = warnings

	return report, nil
}
