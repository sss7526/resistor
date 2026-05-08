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

///////////////////////////////////////////////////////////////////////////////
// 6-Band Deterministic Decode
///////////////////////////////////////////////////////////////////////////////

func TestInferResistor_6BandDeterministic(t *testing.T) {

    // Brown Black Red Brown Brown Brown
    // 1kΩ ±1% 100ppm
    obs := ObservedResistor{
        Bands: []Color{
            Brown, Black, Black, Brown, Brown, Brown,
        },
    }

    res, err := InferResistor(obs)
    require.NoError(t, err)

    require.Equal(t, 1000.0, res.Spec.ResistanceOhms)
    require.Equal(t, 1.0, res.Spec.TolerancePct)
    require.Equal(t, 100, res.Spec.TempCoeffPPM)

    require.InDelta(t, 1.0, res.Meta.Confidence, 1e-9)
}

///////////////////////////////////////////////////////////////////////////////
// Body Color Rules
///////////////////////////////////////////////////////////////////////////////

func TestInferResistor_BodyColorRules(t *testing.T) {

    tests := []struct {
        name     string
        color    Color
        expected ResistorType
    }{
        {"Blue → MetalFilm", Blue, MetalFilm},
        {"Beige → CarbonFilm", Color("beige"), CarbonFilm},
        {"Green → MetalOxide", Green, ResistorType("metal_oxide")},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {

            obs := ObservedResistor{
                BodyColor: tt.color,
            }

            res, err := InferResistor(obs)
            require.NoError(t, err)

            require.Equal(t, tt.expected, res.Spec.Type)
            require.True(t, res.Meta.Confidence > 0)
        })
    }
}

///////////////////////////////////////////////////////////////////////////////
// Length-Based Power Tiers
///////////////////////////////////////////////////////////////////////////////

func TestInferResistor_LengthPowerTiers(t *testing.T) {

    tests := []struct {
        length   float64
        expected float64
    }{
        {2.5, 0.0625},
        {3.5, 0.125},
        {6.0, 0.25},
        {9.0, 0.5},
        {12.0, 1.0},
        {16.0, 2.0},
    }

    for _, tt := range tests {
        t.Run("Length tier", func(t *testing.T) {

            obs := ObservedResistor{
                LengthMM: tt.length,
            }

            res, err := InferResistor(obs)
            require.NoError(t, err)

            require.Equal(t, tt.expected, res.Spec.PowerWatts)
        })
    }
}

///////////////////////////////////////////////////////////////////////////////
// SMD Package Power Rules
///////////////////////////////////////////////////////////////////////////////

func TestInferResistor_SMDPackagePower(t *testing.T) {

    tests := []struct {
        pkg      PackageType
        expected float64
    }{
        {SMD0402, 0.0625},
        {SMD0603, 0.1},
        {SMD0805, 0.125},
        {SMD1206, 0.25},
    }

    for _, tt := range tests {
        t.Run("SMD tier", func(t *testing.T) {

            obs := ObservedResistor{
                Package: tt.pkg,
            }

            res, err := InferResistor(obs)
            require.NoError(t, err)

            require.Equal(t, tt.expected, res.Spec.PowerWatts)
        })
    }
}

///////////////////////////////////////////////////////////////////////////////
// Reinforcement Rule
///////////////////////////////////////////////////////////////////////////////

func TestInferResistor_Reinforcement(t *testing.T) {

    baseObs := ObservedResistor{
        BodyColor: Blue,
    }

    reinforcedObs := ObservedResistor{
        BodyColor: Blue,
        Bands:     []Color{Brown, Black, Red, Brown, Brown}, // 5-band
    }

    baseRes, _ := InferResistor(baseObs)
    reinforcedRes, _ := InferResistor(reinforcedObs)

    require.True(t,
        reinforcedRes.Meta.Confidence >= baseRes.Meta.Confidence)
}
