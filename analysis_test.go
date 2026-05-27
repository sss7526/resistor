package resistor

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAnalyzeResistor_VoltageDriven(t *testing.T) {

	input := AnalysisInput{
		Spec: ResistorSpec{
			ResistanceOhms: 100,
			PowerWatts:     0.5,
			TolerancePct:   5,
		},
		AppliedVoltage: 10,
	}

	report, err := AnalyzeResistor(input)
	require.NoError(t, err)

	require.InDelta(t, 0.1, report.Current, 1e-9)
	require.InDelta(t, 1.0, report.PowerDissipation, 1e-9)
	require.Equal(t, 95.0, report.WorstCaseResistanceMin)
	require.Equal(t, 105.0, report.WorstCaseResistanceMax)

	require.NotEmpty(t, report.Warnings)
}

func TestAnalyzeResistor_CurrentDriven(t *testing.T) {

	input := AnalysisInput{
		Spec: ResistorSpec{
			ResistanceOhms: 50,
		},
		AppliedCurrent: 0.2,
	}

	report, err := AnalyzeResistor(input)
	require.NoError(t, err)

	require.InDelta(t, 10.0, report.VoltageDrop, 1e-9)
	require.InDelta(t, 2.0, report.PowerDissipation, 1e-9)
}

func TestAnalyzeResistor_DeratingWarnings(t *testing.T) {

	input := AnalysisInput{
		Spec: ResistorSpec{
			ResistanceOhms: 100,
			PowerWatts:     0.25,
		},
		AppliedVoltage: 10,
	}

	report, err := AnalyzeResistor(input)
	require.NoError(t, err)

	require.NotEmpty(t, report.Warnings)
}

func TestAnalyzeResistor_OhmsLawConsistency(t *testing.T) {

	// V and I that disagree by more than 1%: V=10, I=0.001, R=100 → expected V=0.1
	inconsistent := AnalysisInput{
		Spec:           ResistorSpec{ResistanceOhms: 100},
		AppliedVoltage: 10,
		AppliedCurrent: 0.001,
	}
	report, err := AnalyzeResistor(inconsistent)
	require.NoError(t, err)

	var found bool
	for _, w := range report.Warnings {
		if w.Level == WarningCaution {
			found = true
			break
		}
	}
	require.True(t, found, "expected a WarningCaution for inconsistent V/I")

	// V and I that agree within 1%: V=10, I=0.1, R=100 → exactly consistent
	consistent := AnalysisInput{
		Spec:           ResistorSpec{ResistanceOhms: 100},
		AppliedVoltage: 10,
		AppliedCurrent: 0.1,
	}
	report2, err := AnalyzeResistor(consistent)
	require.NoError(t, err)

	for _, w := range report2.Warnings {
		require.NotEqual(t, WarningCaution, w.Level, "consistent V/I should not produce a WarningCaution from Ohm's Law check")
	}
}

func TestAnalyzeResistor_NoInputs(t *testing.T) {

	input := AnalysisInput{
		Spec: ResistorSpec{
			ResistanceOhms: 100,
		},
	}

	report, err := AnalyzeResistor(input)
	require.NoError(t, err)

	require.Equal(t, 0.0, report.PowerDissipation)
	require.NotEmpty(t, report.Warnings)
}
