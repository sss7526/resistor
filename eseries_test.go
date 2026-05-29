package resistor

import (
	"testing"

	"github.com/stretchr/testify/require"
)

///////////////////////////////////////////////////////////////////////////////
// RoundNearest Tests
///////////////////////////////////////////////////////////////////////////////

func TestNearestStandard_RoundNearest(t *testing.T) {

	tests := []struct {
		name     string
		input    float64
		series   ESeries
		expected float64
	}{
		{
			name:     "E12 500Ω → 470Ω",
			input:    500,
			series:   E12,
			expected: 470,
		},
		{
			name:     "E24 500Ω → 510Ω",
			input:    500,
			series:   E24,
			expected: 510,
		},
		{
			name:     "E12 4.8kΩ → 4.7kΩ",
			input:    4800,
			series:   E12,
			expected: 4700,
		},
		{
			name:     "E24 4.8kΩ → 4.7kΩ",
			input:    4800,
			series:   E24,
			expected: 4700,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, err := NearestStandard(tt.input, tt.series, RoundNearest)
			require.NoError(t, err)
			require.Equal(t, tt.expected, v)
		})
	}
}

///////////////////////////////////////////////////////////////////////////////
// RoundUp Tests
///////////////////////////////////////////////////////////////////////////////

func TestNearestStandard_RoundUp(t *testing.T) {

	tests := []struct {
		name     string
		input    float64
		series   ESeries
		expected float64
	}{
		{
			name:     "E12 500Ω → 560Ω",
			input:    500,
			series:   E12,
			expected: 560,
		},
		{
			name:     "E24 500Ω → 510Ω",
			input:    500,
			series:   E24,
			expected: 510,
		},
		{
			name:     "E12 4.8kΩ → 5.6kΩ",
			input:    4800,
			series:   E12,
			expected: 5600,
		},
		{
			name:     "E24 9.15Ω → 10Ω (decade boundary)",
			input:    9.15,
			series:   E24,
			expected: 10.0,
		},
		{
			name:     "E24 9150Ω → 10000Ω (decade boundary, kΩ scale)",
			input:    9150,
			series:   E24,
			expected: 10000.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, err := NearestStandard(tt.input, tt.series, RoundUp)
			require.NoError(t, err)
			require.Equal(t, tt.expected, v)
		})
	}
}

///////////////////////////////////////////////////////////////////////////////
// RoundDown Tests
///////////////////////////////////////////////////////////////////////////////

func TestNearestStandard_RoundDown(t *testing.T) {

	tests := []struct {
		name     string
		input    float64
		series   ESeries
		expected float64
	}{
		{
			name:     "E12 500Ω → 470Ω",
			input:    500,
			series:   E12,
			expected: 470,
		},
		{
			name:     "E24 500Ω → 470Ω",
			input:    500,
			series:   E24,
			expected: 470,
		},
		{
			name:     "E12 4.8kΩ → 4.7kΩ",
			input:    4800,
			series:   E12,
			expected: 4700,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, err := NearestStandard(tt.input, tt.series, RoundDown)
			require.NoError(t, err)
			require.Equal(t, tt.expected, v)
		})
	}
}

///////////////////////////////////////////////////////////////////////////////
// Higher Series Sanity Check
///////////////////////////////////////////////////////////////////////////////

func TestNearestStandard_E96_Sanity(t *testing.T) {

	// E96 should snap 499Ω near 499Ω or 499-ish value
	v, err := NearestStandard(499, E96, RoundNearest)
	require.NoError(t, err)

	// Should be within 1% of input (sanity check for generation)
	require.InDelta(t, 499, v, 5)
}

///////////////////////////////////////////////////////////////////////////////
// Branch Coverage
///////////////////////////////////////////////////////////////////////////////

func TestRoundToSignificant_Zero(t *testing.T) {
	require.Equal(t, 0.0, roundToSignificant(0, 3))
}

func TestNearestStandard_RoundNearest_AtFirstBase(t *testing.T) {
	// normalized == base[0] exactly → i=0 branch
	v, err := NearestStandard(1000, E24, RoundNearest)
	require.NoError(t, err)
	require.Equal(t, 1000.0, v)
}

func TestNearestStandard_RoundNearest_CrossDecade(t *testing.T) {
	// normalized=9.8 is closer to 10.0 than to E24's last value (9.1)
	v, err := NearestStandard(9800, E24, RoundNearest)
	require.NoError(t, err)
	require.Equal(t, 10000.0, v)
}

func TestNearestStandard_RoundDown_ExactBaseValue(t *testing.T) {
	// normalized=4.7 lands exactly on a base value → base[i] <= normalized is true
	v, err := NearestStandard(470, E24, RoundDown)
	require.NoError(t, err)
	require.Equal(t, 470.0, v)
}

func TestNearestStandard_RoundNearest_LastBaseBeforeDecode(t *testing.T) {
	// normalized=9.2 > E24 last value (9.1): i==len(base), but 9.2 is closer to 9.1 than to 10.0
	v, err := NearestStandard(9200, E24, RoundNearest)
	require.NoError(t, err)
	require.Equal(t, 9100.0, v)
}

///////////////////////////////////////////////////////////////////////////////
// Error Handling
///////////////////////////////////////////////////////////////////////////////

func TestNearestStandard_Errors(t *testing.T) {

	tests := []struct {
		name   string
		input  float64
		series ESeries
	}{
		{
			name:   "zero value",
			input:  0,
			series: E12,
		},
		{
			name:   "negative value",
			input:  -100,
			series: E12,
		},
		{
			name:   "unsupported series",
			input:  100,
			series: ESeries(999),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NearestStandard(tt.input, tt.series, RoundNearest)
			require.Error(t, err)
		})
	}
}
