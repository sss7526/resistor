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
