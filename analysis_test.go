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