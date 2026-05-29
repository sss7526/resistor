package resistor

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
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
		t.Run(fmt.Sprintf("EIA96 encode %.0f", val), func(t *testing.T) {

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

///////////////////////////////////////////////////////////////////////////////
// E96 Auto-Fallback and Round-Trip Tests
///////////////////////////////////////////////////////////////////////////////

// TestSMD_E96_AutoFallback verifies that SMDAuto falls through to EIA-96 for
// E96 values that are not representable in 3/4-digit format.
func TestSMD_E96_AutoFallback(t *testing.T) {
	cases := []struct {
		name  string
		value float64
	}{
		{"1.02Ω", 1.02},
		{"1.05Ω", 1.05},
		{"10.2Ω", 10.2},
		{"10.5Ω", 10.5},
		{"52.3Ω", 52.3},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			marking, err := EncodeSMD(tt.value, SMDAuto)
			require.NoError(t, err)
			spec, err := DecodeSMD(marking)
			require.NoError(t, err)
			require.InDelta(t, tt.value, spec.ResistanceOhms, 1e-6)
		})
	}
}

// TestSMD_E96_RoundTrip_AllBaseValues tests DecodeSMD(EncodeSMD(v, SMDAuto)) == v
// for all 96 E96 base values across all 12 EIA-96 multiplier decades.
func TestSMD_E96_RoundTrip_AllBaseValues(t *testing.T) {
	for letter, multiplier := range eia96Multipliers {
		for i, base := range eia96Base {
			v := roundToSignificant(base*multiplier, 6)
			t.Run(fmt.Sprintf("base[%02d]x%c=%.6g", i+1, letter, v), func(t *testing.T) {
				marking, err := EncodeSMD(v, SMDAuto)
				require.NoError(t, err, "EncodeSMD failed")
				spec, err := DecodeSMD(marking)
				require.NoError(t, err, "DecodeSMD failed on marking %q", marking)
				require.InDelta(t, v, spec.ResistanceOhms, v*1e-5,
					"round-trip mismatch: want %.6g, got %.6g via %q", v, spec.ResistanceOhms, marking)
			})
		}
	}
}

// TestSMD_E96_NearestStandard_RoundTrip verifies that NearestStandard(E96) results
// feed into EncodeSMD without error for all 96 base values across common decades.
func TestSMD_E96_NearestStandard_RoundTrip(t *testing.T) {
	multipliers := []float64{10, 100, 1_000, 10_000, 100_000}
	for _, mult := range multipliers {
		for _, base := range eia96Base {
			v := roundToSignificant(base*mult, 6)
			snapped, err := NearestStandard(v, E96, RoundNearest)
			require.NoError(t, err)
			require.InDelta(t, v, snapped, v*1e-5, "NearestStandard should return exact E96 value")
			_, err = EncodeSMD(snapped, SMDAuto)
			require.NoError(t, err, "EncodeSMD(SMDAuto) failed for E96 value %.6g", snapped)
		}
	}
}

// TestSMD_E96_SMDStandard_StillStrict verifies that SMDStandard still rejects
// values that are not representable in 3/4-digit format.
func TestSMD_E96_SMDStandard_StillStrict(t *testing.T) {
	nonStandardValues := []float64{1.02, 10.5, 52.3}
	for _, v := range nonStandardValues {
		t.Run(fmt.Sprintf("%.4g", v), func(t *testing.T) {
			_, err := EncodeSMD(v, SMDStandard)
			require.Error(t, err, "SMDStandard should still reject %.4g", v)
		})
	}
}

// TestSMD_E96_AutoNonRepresentable verifies that SMDAuto still errors when a
// value is neither 3/4-digit nor EIA-96 representable.
func TestSMD_E96_AutoNonRepresentable(t *testing.T) {
	// 1234 is not in the E96 series and not 3/4-digit representable
	// (1234/10=123.4, not exact integer; 1234/100=12.34, not exact)
	_, err := EncodeSMD(1234.5678, SMDAuto)
	require.Error(t, err)
}

func TestSMD_Encode_UnsupportedMode(t *testing.T) {
	_, err := EncodeSMD(1000, SMDEncodingMode(99))
	require.Error(t, err)
}

func TestSMD_Decode_BareR_ZeroOhm(t *testing.T) {
	// "R" decodes as 0.0 Ω via R-notation (left="" right=""), which is invalid
	_, err := DecodeSMD("R")
	require.Error(t, err)
}

func TestSMD_Decode_RNotation_ParseFloatError(t *testing.T) {
	// "AR7" passes the R-notation check but "A.7" fails strconv.ParseFloat
	_, err := DecodeSMD("AR7")
	require.Error(t, err)
}
