package resistor

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBodyColors(t *testing.T) {
	colors := BodyColors()
	require.Contains(t, colors, Blue)
	require.Contains(t, colors, Green)
	require.Contains(t, colors, Color("beige"))
	require.Contains(t, colors, Color("tan"))
}
