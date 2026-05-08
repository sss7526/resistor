package resistor

import (
	"github.com/stretchr/testify/require"
	"testing"
)

///////////////////////////////////////////////////////////////////////////////
// Decode Tests
///////////////////////////////////////////////////////////////////////////////

func TestSMD_Decode_ValidCases(t *testing.T) {

	tests := []struct {
		name     string
		marking  string
		expected float64
	}{
		{
			name:     "3-digit 472",
			marking:  "472",
			expected: 4700,
		},
		{
			name:     "4-digit 4701",
			marking:  "4701",
			expected: 4700,
		},
		{
			name:     "R notation 4R7",
			marking:  "4R7",
			expected: 4.7,
		},
		{
			name:     "R notation R47",
			marking:  "R47",
			expected: 0.47,
		},
		{
			name:     "lowercase r notation",
			marking:  "4r7",
			expected: 4.7,
		},
		{
			name:     "EIA96 01X",
			marking:  "01X",
			expected: roundToSignificant(eia96Base[0]*1, 6),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec, err := DecodeSMD(tt.marking)
			require.NoError(t, err)
			require.InDelta(t, tt.expected, spec.ResistanceOhms, 1e-6)
		})
	}
}

func TestSMD_Decode_InvalidCases(t *testing.T) {

	tests := []string{
		"",
		"ABC",
		"12",
		"99999",
		"01Q",  // invalid EIA multiplier
		"97X",  // invalid EIA code
		"4RR7", // invalid R notation
	}

	for _, m := range tests {
		t.Run(m, func(t *testing.T) {
			_, err := DecodeSMD(m)
			require.Error(t, err)
		})
	}
}

///////////////////////////////////////////////////////////////////////////////
// Encode Standard Tests
///////////////////////////////////////////////////////////////////////////////

func TestSMD_Encode_Standard(t *testing.T) {

	tests := []struct {
		name     string
		value    float64
		expected string
	}{
		{
			name:     "4700Ω",
			value:    4700,
			expected: "472",
		},
		{
			name:     "100Ω",
			value:    100,
			expected: "101",
		},
		{
			name:     "1kΩ",
			value:    1000,
			expected: "102",
		},
		{
			name:     "4.7kΩ 4-digit",
			value:    4700,
			expected: "472",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, err := EncodeSMD(tt.value, SMDStandard)
			require.NoError(t, err)
			require.Equal(t, tt.expected, code)
		})
	}
}

///////////////////////////////////////////////////////////////////////////////
// Encode EIA-96 Tests
///////////////////////////////////////////////////////////////////////////////

func TestSMD_Encode_EIA96(t *testing.T) {

	tests := []float64{
		100,  // 100Ω
		1000, // 1kΩ
		4990, // typical E96 value
	}

	for _, val := range tests {
		t.Run("EIA96 encode "+string(rune(int(val))), func(t *testing.T) {

			code, err := EncodeSMD(val, SMDEIA96)
			require.NoError(t, err)

			spec, err := DecodeSMD(code)
			require.NoError(t, err)

			require.InDelta(t, val, spec.ResistanceOhms, 1e-6)
		})
	}
}

///////////////////////////////////////////////////////////////////////////////
// Auto Mode Tests
///////////////////////////////////////////////////////////////////////////////

func TestSMD_Encode_AutoMode(t *testing.T) {

	code, err := EncodeSMD(4700, SMDAuto)
	require.NoError(t, err)
	require.Equal(t, "472", code)
}

///////////////////////////////////////////////////////////////////////////////
// Encode Invalid Cases
///////////////////////////////////////////////////////////////////////////////

func TestSMD_Encode_InvalidCases(t *testing.T) {

	tests := []struct {
		name  string
		value float64
		mode  SMDEncodingMode
	}{
		{
			name:  "zero value",
			value: 0,
			mode:  SMDStandard,
		},
		{
			name:  "negative value",
			value: -100,
			mode:  SMDStandard,
		},
		{
			name:  "non representable EIA96",
			value: 1234,
			mode:  SMDEIA96,
		},
		{
			name:  "non representable standard",
			value: 123.456,
			mode:  SMDStandard,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := EncodeSMD(tt.value, tt.mode)
			require.Error(t, err)
		})
	}
}

///////////////////////////////////////////////////////////////////////////////
// Round Trip Standard
///////////////////////////////////////////////////////////////////////////////

func TestSMD_RoundTrip_Standard(t *testing.T) {

	values := []float64{
		100,
		470,
		1000,
		4700,
	}

	for _, v := range values {
		t.Run("roundtrip", func(t *testing.T) {

			code, err := EncodeSMD(v, SMDStandard)
			require.NoError(t, err)

			spec, err := DecodeSMD(code)
			require.NoError(t, err)

			require.InDelta(t, v, spec.ResistanceOhms, 1e-6)
		})
	}
}

///////////////////////////////////////////////////////////////////////////////
// Round Trip EIA-96
///////////////////////////////////////////////////////////////////////////////

func TestSMD_RoundTrip_EIA96(t *testing.T) {

	values := []float64{
		100,
		1000,
		4990,
	}

	for _, v := range values {
		t.Run("roundtrip EIA96", func(t *testing.T) {

			code, err := EncodeSMD(v, SMDEIA96)
			require.NoError(t, err)

			spec, err := DecodeSMD(code)
			require.NoError(t, err)

			require.InDelta(t, v, spec.ResistanceOhms, 1e-6)
		})
	}
}
