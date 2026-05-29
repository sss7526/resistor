package resistor

import (
	"testing"

	"github.com/stretchr/testify/require"
)

///////////////////////////////////////////////////////////////////////////////
// ESeries
///////////////////////////////////////////////////////////////////////////////

func TestESeries_String(t *testing.T) {
	tests := []struct {
		series ESeries
		want   string
	}{
		{E3, "E3"},
		{E6, "E6"},
		{E12, "E12"},
		{E24, "E24"},
		{E48, "E48"},
		{E96, "E96"},
		{E192, "E192"},
		{ESeries(999), "Unknown"},
	}
	for _, tt := range tests {
		require.Equal(t, tt.want, tt.series.String())
	}
}

func TestParseESeries_Valid(t *testing.T) {
	tests := []struct {
		input string
		want  ESeries
	}{
		{"E3", E3},
		{"E6", E6},
		{"E12", E12},
		{"E24", E24},
		{"E48", E48},
		{"E96", E96},
		{"E192", E192},
		{"e24", E24},   // case-insensitive
		{" E24 ", E24}, // trimmed
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParseESeries(tt.input)
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestParseESeries_Invalid(t *testing.T) {
	for _, input := range []string{"", "E5", "E25", "24", "invalid"} {
		t.Run(input, func(t *testing.T) {
			_, err := ParseESeries(input)
			require.Error(t, err)
		})
	}
}

func TestAllESeries(t *testing.T) {
	all := AllESeries()
	require.Len(t, all, 7)
	require.Contains(t, all, E3)
	require.Contains(t, all, E6)
	require.Contains(t, all, E12)
	require.Contains(t, all, E24)
	require.Contains(t, all, E48)
	require.Contains(t, all, E96)
	require.Contains(t, all, E192)
}

///////////////////////////////////////////////////////////////////////////////
// RoundingMode
///////////////////////////////////////////////////////////////////////////////

func TestRoundingMode_String(t *testing.T) {
	tests := []struct {
		mode RoundingMode
		want string
	}{
		{RoundNearest, "Nearest"},
		{RoundUp, "Up"},
		{RoundDown, "Down"},
		{RoundingMode(99), "Unknown"},
	}
	for _, tt := range tests {
		require.Equal(t, tt.want, tt.mode.String())
	}
}

func TestParseRoundingMode_Valid(t *testing.T) {
	tests := []struct {
		input string
		want  RoundingMode
	}{
		{"nearest", RoundNearest},
		{"Nearest", RoundNearest},
		{"", RoundNearest}, // empty defaults to nearest
		{"up", RoundUp},
		{"Up", RoundUp},
		{"down", RoundDown},
		{"Down", RoundDown},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParseRoundingMode(tt.input)
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestParseRoundingMode_Invalid(t *testing.T) {
	for _, input := range []string{"center", "random", "nearest_up"} {
		t.Run(input, func(t *testing.T) {
			_, err := ParseRoundingMode(input)
			require.Error(t, err)
		})
	}
}

func TestAllRoundingModes(t *testing.T) {
	all := AllRoundingModes()
	require.Len(t, all, 3)
	require.Contains(t, all, RoundNearest)
	require.Contains(t, all, RoundUp)
	require.Contains(t, all, RoundDown)
}

///////////////////////////////////////////////////////////////////////////////
// PackageType
///////////////////////////////////////////////////////////////////////////////

func TestPackageType_String(t *testing.T) {
	tests := []struct {
		pkg  PackageType
		want string
	}{
		{ThroughHole, "through_hole"},
		{SMD, "smd"},
		{SMD0402, "0402"},
		{SMD0603, "0603"},
		{SMD0805, "0805"},
		{SMD1206, "1206"},
		{Axial, "axial"},
		{Radial, "radial"},
		{UnknownPKG, "unknown"},
		{PackageType("bogus"), "unknown"},
	}
	for _, tt := range tests {
		require.Equal(t, tt.want, tt.pkg.String())
	}
}

func TestParsePackageType_Valid(t *testing.T) {
	tests := []struct {
		input string
		want  PackageType
	}{
		{"through_hole", ThroughHole},
		{"Through_Hole", ThroughHole}, // case-insensitive
		{"smd", SMD},
		{"0402", SMD0402},
		{"0603", SMD0603},
		{"0805", SMD0805},
		{"1206", SMD1206},
		{"axial", Axial},
		{"radial", Radial},
		{"", PackageType("")}, // empty → zero value, no error
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParsePackageType(tt.input)
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestParsePackageType_Invalid(t *testing.T) {
	for _, input := range []string{"unknown", "0201", "2512", "dip"} {
		t.Run(input, func(t *testing.T) {
			_, err := ParsePackageType(input)
			require.Error(t, err)
		})
	}
}

func TestAllPackageTypes(t *testing.T) {
	all := AllPackageTypes()
	require.NotEmpty(t, all)
	require.Contains(t, all, ThroughHole)
	require.Contains(t, all, SMD0402)
	require.Contains(t, all, SMD0603)
	require.Contains(t, all, SMD0805)
	require.Contains(t, all, SMD1206)
	require.Contains(t, all, Axial)
	require.Contains(t, all, Radial)
	for _, pt := range all {
		require.NotEqual(t, UnknownPKG, pt, "AllPackageTypes should not include UnknownPKG")
	}
}

///////////////////////////////////////////////////////////////////////////////
// Color Enumerations
///////////////////////////////////////////////////////////////////////////////

func TestDigitColors(t *testing.T) {
	colors := DigitColors()
	require.Len(t, colors, 10)
	require.Contains(t, colors, Black)
	require.Contains(t, colors, White)
}

func TestMultiplierColors(t *testing.T) {
	colors := MultiplierColors()
	require.Len(t, colors, 12)
	require.Contains(t, colors, Gold)
	require.Contains(t, colors, Silver)
}

func TestToleranceColors(t *testing.T) {
	colors := ToleranceColors()
	require.Len(t, colors, 9)
	require.Contains(t, colors, Gold)
	require.Contains(t, colors, None)
}

func TestTempCoeffColors(t *testing.T) {
	colors := TempCoeffColors()
	require.Len(t, colors, 6)
	require.Contains(t, colors, Brown)
	require.Contains(t, colors, Violet)
}
