package resistor

import (
	"fmt"
	"math"
)

///////////////////////////////////////////////////////////////////////////////
// Warning System
///////////////////////////////////////////////////////////////////////////////

// WarningLevel indicates the severity of an analysis warning.
type WarningLevel string

const (
	// WarningInfo is an informational note; no action required.
	WarningInfo WarningLevel = "info"

	// WarningCaution indicates a condition that may warrant review.
	WarningCaution WarningLevel = "caution"

	// WarningDanger indicates a condition likely to cause component failure.
	WarningDanger WarningLevel = "danger"
)

// AnalysisWarning represents a structured engineering warning produced during analysis.
type AnalysisWarning struct {
	// Level is the severity classification of the warning.
	Level WarningLevel

	// Message is a human-readable description of the condition.
	Message string
}

///////////////////////////////////////////////////////////////////////////////
// Analysis Models
///////////////////////////////////////////////////////////////////////////////

// AnalysisInput contains electrical conditions for resistor analysis.
//
// Either AppliedVoltage or AppliedCurrent (or both) may be provided.
// If both are provided, consistency is checked against Ohm's Law.
type AnalysisInput struct {
	// Spec is the resistor being analyzed.
	Spec ResistorSpec

	// AppliedVoltage is the voltage across the resistor in volts.
	// Leave zero if not known.
	AppliedVoltage float64

	// AppliedCurrent is the current through the resistor in amperes.
	// Leave zero if not known.
	AppliedCurrent float64
}

// AnalysisReport contains deterministic electrical analysis results.
//
// Pointer fields (DeratedSafePower, WorstCaseResistanceMin, WorstCaseResistanceMax)
// are nil when the required inputs were absent. A non-nil pointer to zero is a
// legitimately computed value (e.g. WorstCaseResistanceMin at 100% tolerance).
type AnalysisReport struct {
	// PowerDissipation is the computed power dissipation in watts (P = I²R).
	PowerDissipation float64

	// VoltageDrop is the voltage across the resistor in volts.
	VoltageDrop float64

	// Current is the current through the resistor in amperes.
	Current float64

	// DeratedSafePower is the 50%-derated safe operating power in watts.
	// Nil if the rated power was not provided.
	DeratedSafePower *float64 `json:"DeratedSafePower,omitempty"`

	// WorstCaseResistanceMin is the minimum resistance at the tolerance limit in ohms.
	// Nil if tolerance was not provided.
	WorstCaseResistanceMin *float64 `json:"WorstCaseResistanceMin,omitempty"`

	// WorstCaseResistanceMax is the maximum resistance at the tolerance limit in ohms.
	// Nil if tolerance was not provided.
	WorstCaseResistanceMax *float64 `json:"WorstCaseResistanceMax,omitempty"`

	// Warnings contains any engineering warnings generated during analysis.
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
