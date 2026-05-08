package resistor

import (
	"testing"
)

func FuzzDecodeBands(f *testing.F) {

	// Seed with valid known examples
	f.Add("brown,black,red,gold")
	f.Add("brown,black.black,brown,brown")
	f.Add("brown,black,black,brown,brown,brown")

	f.Fuzz(func(t *testing.T, input string) {

		parts := splitAndTrim(input)
		var bands []Color

		for _, p := range parts {
			bands = append(bands, Color(p))
		}

		_, _ = DecodeBands(bands)
	})
}

func splitAndTrim(s string) []string {
	var out []string
	start := 0
	for i := 0; i <= len(s); i++ {
		if i == len(s) || s[i] == ',' {
			if start < i {
				out = append(out, s[start:i])
			}
			start = i + 1
		}
	}
	return out
}

func FuzzDecodeSMD(f *testing.F) {

	// Seed valid examples
	f.Add("472")
	f.Add("4701")
	f.Add("4R7")
	f.Add("01C")

	f.Fuzz(func(t *testing.T, input string) {
		_, _ = DecodeSMD(input)
	})
}

func FuzzNearestStandard(f *testing.F) {

	f.Add(100.0)
	f.Add(4700.0)
	f.Add(1.0)

	f.Fuzz(func(t *testing.T, val float64) {

		if val <= 0 || val > 1e9 {
			return
		}

		_, _ = NearestStandard(val, E24, RoundNearest)
	})
}
