package resistor

import (
    "testing"

    "github.com/stretchr/testify/require"
)

///////////////////////////////////////////////////////////////////////////////
// Deterministic Only (Color Bands)
///////////////////////////////////////////////////////////////////////////////

func TestInferResistor_DeterministicBands(t *testing.T) {

    obs := ObservedResistor{
        Bands: []Color{Yellow, Violet, Brown, Gold}, // 470Ω ±5%
    }

    res, err := InferResistor(obs)
    require.NoError(t, err)

    require.Equal(t, 470.0, res.Spec.ResistanceOhms)
    require.Equal(t, 5.0, res.Spec.TolerancePct)

    require.Contains(t, res.Meta.Assumptions,
        "Resistance and tolerance determined from color bands")

    require.InDelta(t, 1.0, res.Meta.Confidence, 1e-9)
}

///////////////////////////////////////////////////////////////////////////////
// Deterministic Only (SMD)
///////////////////////////////////////////////////////////////////////////////

func TestInferResistor_DeterministicSMD(t *testing.T) {

    obs := ObservedResistor{
        Marking: "472",
    }

    res, err := InferResistor(obs)
    require.NoError(t, err)

    require.Equal(t, 4700.0, res.Spec.ResistanceOhms)

    require.Contains(t, res.Meta.Assumptions,
        "Resistance determined from SMD marking")

    require.InDelta(t, 1.0, res.Meta.Confidence, 1e-9)
}

///////////////////////////////////////////////////////////////////////////////
// Heuristic Only (No Deterministic Resistance)
///////////////////////////////////////////////////////////////////////////////

func TestInferResistor_HeuristicOnly(t *testing.T) {

    obs := ObservedResistor{
        BodyColor: Blue,
        LengthMM:  6.2,
    }

    res, err := InferResistor(obs)
    require.NoError(t, err)

    require.Equal(t, MetalFilm, res.Spec.Type)
    require.Equal(t, 0.25, res.Spec.PowerWatts)

    require.NotEmpty(t, res.Meta.Assumptions)

    require.True(t, res.Meta.Confidence > 0)
    require.True(t, res.Meta.Confidence < 1)
}

///////////////////////////////////////////////////////////////////////////////
// Combined Deterministic + Heuristic
///////////////////////////////////////////////////////////////////////////////

func TestInferResistor_Combined(t *testing.T) {

    obs := ObservedResistor{
        Bands:     []Color{Brown, Black, Red, Gold}, // 1kΩ ±5%
        BodyColor: Blue,
        LengthMM:  6.3,
    }

    res, err := InferResistor(obs)
    require.NoError(t, err)

    require.Equal(t, 1000.0, res.Spec.ResistanceOhms)
    require.Equal(t, 5.0, res.Spec.TolerancePct)
    require.Equal(t, MetalFilm, res.Spec.Type)
    require.Equal(t, 0.25, res.Spec.PowerWatts)

    require.True(t, res.Meta.Confidence > 0.7)
    require.True(t, res.Meta.Confidence <= 1.0)
}

///////////////////////////////////////////////////////////////////////////////
// Deterministic Overrides Heuristic
///////////////////////////////////////////////////////////////////////////////

func TestInferResistor_DeterministicOverrides(t *testing.T) {

    obs := ObservedResistor{
        Bands: []Color{Brown, Black, Red, Brown, Brown}, // 1kΩ ±1%
    }

    res, err := InferResistor(obs)
    require.NoError(t, err)

    require.Equal(t, 1.0, res.Spec.TolerancePct)

    // 5-band heuristic should not override deterministic tolerance
    require.NotContains(t, res.Meta.Assumptions,
        "5 bands assumed ±1% tolerance")
}

///////////////////////////////////////////////////////////////////////////////
// Confidence Monotonicity
///////////////////////////////////////////////////////////////////////////////

func TestInferResistor_ConfidenceMonotonicity(t *testing.T) {

    baseObs := ObservedResistor{
        BodyColor: Blue,
    }

    extendedObs := ObservedResistor{
        BodyColor: Blue,
        LengthMM:  6.3,
    }

    baseRes, _ := InferResistor(baseObs)
    extendedRes, _ := InferResistor(extendedObs)

    require.True(t,
        extendedRes.Meta.Confidence >= baseRes.Meta.Confidence)
}

///////////////////////////////////////////////////////////////////////////////
// Empty Observation
///////////////////////////////////////////////////////////////////////////////

func TestInferResistor_EmptyObservation(t *testing.T) {

    obs := ObservedResistor{}

    res, err := InferResistor(obs)
    require.NoError(t, err)

    require.Equal(t, 0.0, res.Spec.ResistanceOhms)
    require.Equal(t, 0.0, res.Meta.Confidence)
    require.Empty(t, res.Meta.Assumptions)
}