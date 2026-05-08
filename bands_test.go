package resistor

import (
	"github.com/stretchr/testify/require"
	"testing"
)

///////////////////////////////////////////////////////////////////////////////
// Decode Tests
///////////////////////////////////////////////////////////////////////////////

func TestDecodeBands_ValidCases(t *testing.T) {

	tests := []struct {
		name       string
		bands      []Color
		wantOhms   float64
		wantTolPct float64
	}{
		{
			name:       "4-band 510Ω ±5%",
			bands:      []Color{Green, Brown, Brown, Gold},
			wantOhms:   510,
			wantTolPct: 5,
		},
		{
			name:       "4-band 1kΩ ±10%",
			bands:      []Color{Brown, Black, Red, Silver},
			wantOhms:   1000,
			wantTolPct: 10,
		},
		{
			name:       "5-band 4.7kΩ ±1%",
			bands:      []Color{Yellow, Violet, Black, Brown, Brown},
			wantOhms:   4700,
			wantTolPct: 1,
		},
		{
			name:       "5-band 100Ω ±2%",
			bands:      []Color{Brown, Black, Black, Black, Red},
			wantOhms:   100,
			wantTolPct: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec, err := DecodeBands(tt.bands)
			require.NoError(t, err)
			require.Equal(t, tt.wantOhms, spec.ResistanceOhms)
			require.Equal(t, tt.wantTolPct, spec.TolerancePct)
		})
	}
}

func TestDecodeBands_InvalidCases(t *testing.T) {

	tests := []struct {
		name  string
		bands []Color
	}{
		{
			name:  "invalid band count",
			bands: []Color{Brown, Black, Red},
		},
		{
			name:  "invalid digit color",
			bands: []Color{Gold, Black, Red, Gold},
		},
		{
			name:  "invalid multiplier",
			bands: []Color{Brown, Black, None, Gold},
		},
		{
			name:  "invalid tolerance",
			bands: []Color{Brown, Black, Red, Black},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := DecodeBands(tt.bands)
			require.Error(t, err)
		})
	}
}

///////////////////////////////////////////////////////////////////////////////
// Encode Tests
///////////////////////////////////////////////////////////////////////////////

func TestEncodeBands_ValidCases(t *testing.T) {

	tests := []struct {
		name       string
		resistance float64
		tolerance  float64
		wantBands  []Color
	}{
		{
			name:       "4-band 510Ω ±5%",
			resistance: 510,
			tolerance:  5,
			wantBands:  []Color{Green, Brown, Brown, Gold},
		},
		{
			name:       "4-band 1kΩ ±10%",
			resistance: 1000,
			tolerance:  10,
			wantBands:  []Color{Brown, Black, Red, Silver},
		},
		{
			name:       "5-band 4.7kΩ ±1%",
			resistance: 4700,
			tolerance:  1,
			wantBands:  []Color{Yellow, Violet, Black, Brown, Brown},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bands, err := EncodeBandsSimple(tt.resistance, tt.tolerance)
			require.NoError(t, err)
			require.Equal(t, tt.wantBands, bands)
		})
	}
}

func TestEncodeBands_InvalidCases(t *testing.T) {

	tests := []struct {
		name       string
		resistance float64
		tolerance  float64
	}{
		{
			name:       "zero resistance",
			resistance: 0,
			tolerance:  5,
		},
		{
			name:       "non-encodable fractional",
			resistance: 123.45,
			tolerance:  5,
		},
		{
			name:       "unsupported tolerance",
			resistance: 1000,
			tolerance:  3, // not a standard tolerance band
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := EncodeBandsSimple(tt.resistance, tt.tolerance)
			require.Error(t, err)
		})
	}
}

///////////////////////////////////////////////////////////////////////////////
// Round-Trip Tests
///////////////////////////////////////////////////////////////////////////////

func TestEncodeDecode_RoundTrip(t *testing.T) {

	tests := []struct {
		name       string
		resistance float64
		tolerance  float64
	}{
		{
			name:       "510Ω ±5%",
			resistance: 510,
			tolerance:  5,
		},
		{
			name:       "1kΩ ±10%",
			resistance: 1000,
			tolerance:  10,
		},
		{
			name:       "4.7kΩ ±1%",
			resistance: 4700,
			tolerance:  1,
		},
		{
			name:       "100Ω ±2%",
			resistance: 100,
			tolerance:  2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			bands, err := EncodeBandsSimple(tt.resistance, tt.tolerance)
			require.NoError(t, err)

			spec, err := DecodeBands(bands)
			require.NoError(t, err)

			require.Equal(t, tt.resistance, spec.ResistanceOhms)
			require.Equal(t, tt.tolerance, spec.TolerancePct)
		})
	}
}

func TestEncodeBands_SixBand(t *testing.T) {

	spec := ResistorSpec{
		ResistanceOhms: 1000,
		TolerancePct:   1,
		TempCoeffPPM:   100,
	}

	bands, err := EncodeBands(spec)
	require.NoError(t, err)

	expected := []Color{
		Brown, Black, Black, Brown, Brown, Brown,
	}

	require.Equal(t, expected, bands)
}
