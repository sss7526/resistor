package resistor

import (
	"testing"

	"github.com/stretchr/testify/require"
)

///////////////////////////////////////////////////////////////////////////////
// BandRole.String
///////////////////////////////////////////////////////////////////////////////

func TestBandRole_String(t *testing.T) {
	tests := []struct {
		role BandRole
		want string
	}{
		{RoleDigit, "Digit"},
		{RoleMultiplier, "Multiplier"},
		{RoleTolerance, "Tolerance"},
		{RoleTempCoeff, "TempCoeff"},
		{BandRole(99), "Unknown"},
	}
	for _, tt := range tests {
		require.Equal(t, tt.want, tt.role.String())
	}
}

///////////////////////////////////////////////////////////////////////////////
// BandRolesForCount
///////////////////////////////////////////////////////////////////////////////

func TestBandRolesForCount_Valid(t *testing.T) {
	tests := []struct {
		count int
		want  []BandRole
	}{
		{
			count: 4,
			want:  []BandRole{RoleDigit, RoleDigit, RoleMultiplier, RoleTolerance},
		},
		{
			count: 5,
			want:  []BandRole{RoleDigit, RoleDigit, RoleDigit, RoleMultiplier, RoleTolerance},
		},
		{
			count: 6,
			want:  []BandRole{RoleDigit, RoleDigit, RoleDigit, RoleMultiplier, RoleTolerance, RoleTempCoeff},
		},
	}
	for _, tt := range tests {
		t.Run(tt.want[0].String(), func(t *testing.T) {
			roles, err := BandRolesForCount(tt.count)
			require.NoError(t, err)
			require.Equal(t, tt.want, roles)
		})
	}
}

func TestBandRolesForCount_Invalid(t *testing.T) {
	for _, count := range []int{0, 1, 2, 3, 7, 100} {
		t.Run("", func(t *testing.T) {
			_, err := BandRolesForCount(count)
			require.Error(t, err)
		})
	}
}

///////////////////////////////////////////////////////////////////////////////
// ValidColorsForRole
///////////////////////////////////////////////////////////////////////////////

func TestValidColorsForRole(t *testing.T) {
	tests := []struct {
		role BandRole
		want []Color
	}{
		{RoleDigit, DigitColors()},
		{RoleMultiplier, MultiplierColors()},
		{RoleTolerance, ToleranceColors()},
		{RoleTempCoeff, TempCoeffColors()},
	}
	for _, tt := range tests {
		t.Run(tt.role.String(), func(t *testing.T) {
			got := ValidColorsForRole(tt.role)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestValidColorsForRole_Unknown(t *testing.T) {
	require.Nil(t, ValidColorsForRole(BandRole(99)))
}
