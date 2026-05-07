package resistor

import (
	"testing"

	"github.com/stretchr/testify/require"
)

///////////////////////////////////////////////////////////////////////////////
// Basic Default Behavior
///////////////////////////////////////////////////////////////////////////////

func TestSelectStandardResistor_DefaultsApplied(t *testing.T) {

	req := SelectionRequest{
		Resistance: 487,
	}

	res, err := SelectStandardResistor(req)
	require.NoError(t, err)

	require.Equal(t, 487.0, res.RequestedResistance)
	require.Equal(t, 470.0, res.SelectedResistance)
	require.Equal(t, E24, res.Series)
	require.Equal(t, 5.0, res.TolerancePct)
	require.Equal(t, RoundNearest, res.Rounding)

	require.NotEmpty(t, res.Bands)
	require.Contains(t, res.Assumptions, "Series defaulted to E24")
	require.Contains(t, res.Assumptions, "Tolerance defaulted to ±5%")
	require.Contains(t, res.Assumptions, "Rounding mode defaulted to RoundNearest")
	require.Len(t, res.Assumptions, 4)
}

///////////////////////////////////////////////////////////////////////////////
// Exact Match (No Snap)
///////////////////////////////////////////////////////////////////////////////

func TestSelectStandardResistor_ExactMatch(t *testing.T) {

	req := SelectionRequest{
		Resistance:   470,
		Series:       E24,
		TolerancePct: 5,
		Rounding:     RoundNearest,
	}

	res, err := SelectStandardResistor(req)
	require.NoError(t, err)

	require.Equal(t, 470.0, res.SelectedResistance)
	require.Empty(t, res.Assumptions)
}

///////////////////////////////////////////////////////////////////////////////
// Custom Series Behavior
///////////////////////////////////////////////////////////////////////////////

func TestSelectStandardResistor_CustomSeries(t *testing.T) {

	req := SelectionRequest{
		Resistance: 500,
		Series:     E12,
	}

	res, err := SelectStandardResistor(req)
	require.NoError(t, err)

	require.Equal(t, E12, res.Series)
	require.Equal(t, 470.0, res.SelectedResistance)
	require.Contains(t, res.Assumptions, "Tolerance defaulted to ±5%")
	require.Contains(t, res.Assumptions, "Rounding mode defaulted to RoundNearest")
}

///////////////////////////////////////////////////////////////////////////////
// Custom Rounding Behavior
///////////////////////////////////////////////////////////////////////////////

func TestSelectStandardResistor_RoundUp(t *testing.T) {

	req := SelectionRequest{
		Resistance: 500,
		Series:     E12,
		Rounding:   RoundUp,
	}

	res, err := SelectStandardResistor(req)
	require.NoError(t, err)

	require.Equal(t, 560.0, res.SelectedResistance)
	require.Contains(t, res.Assumptions,
		"Resistance snapped from 500Ω to 560Ω",
	)
}

///////////////////////////////////////////////////////////////////////////////
// Custom Tolerance Behavior
///////////////////////////////////////////////////////////////////////////////

func TestSelectStandardResistor_CustomTolerance(t *testing.T) {

	req := SelectionRequest{
		Resistance:   4700,
		TolerancePct: 1,
	}

	res, err := SelectStandardResistor(req)
	require.NoError(t, err)

	require.Equal(t, 1.0, res.TolerancePct)
	require.Len(t, res.Bands, 5) // 5-band expected for ≤2%
}

///////////////////////////////////////////////////////////////////////////////
// Integration: Verify Band Encoding
///////////////////////////////////////////////////////////////////////////////

func TestSelectStandardResistor_BandsCorrect(t *testing.T) {

	req := SelectionRequest{
		Resistance: 500,
	}

	res, err := SelectStandardResistor(req)
	require.NoError(t, err)

	expectedBands := []Color{Green, Brown, Brown, Gold}
	require.Equal(t, expectedBands, res.Bands)
}

///////////////////////////////////////////////////////////////////////////////
// Error Cases
///////////////////////////////////////////////////////////////////////////////

func TestSelectStandardResistor_InvalidInput(t *testing.T) {

	tests := []struct {
		name string
		req  SelectionRequest
	}{
		{
			name: "zero resistance",
			req:  SelectionRequest{Resistance: 0},
		},
		{
			name: "negative resistance",
			req:  SelectionRequest{Resistance: -100},
		},
		{
			name: "unsupported series",
			req:  SelectionRequest{Resistance: 100, Series: ESeries(999)},
		},
		{
			name: "unsupported tolerance",
			req:  SelectionRequest{Resistance: 100, TolerancePct: 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := SelectStandardResistor(tt.req)
			require.Error(t, err)
		})
	}
}
