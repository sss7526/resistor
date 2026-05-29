package resistor

import (
	"errors"
	"math"
	"sort"
)

///////////////////////////////////////////////////////////////////////////////
// E-Series Definition
///////////////////////////////////////////////////////////////////////////////

/*
ESeries represents a standard IEC 60063 preferred number series.

The number indicates how many logarithmically spaced values
exist per decade (between 1.0 and 10.0).

Examples:

E6   → 6 values per decade
E12  → 12 values per decade
E24  → 24 values per decade
E96  → 96 values per decade
E192 → 192 values per decade

These series are used to standardize commercially available resistor values.
*/
type ESeries int

const (
	E3   ESeries = 3
	E6   ESeries = 6
	E12  ESeries = 12
	E24  ESeries = 24
	E48  ESeries = 48
	E96  ESeries = 96
	E192 ESeries = 192
)

///////////////////////////////////////////////////////////////////////////////
// Rounding Modes
///////////////////////////////////////////////////////////////////////////////

/*
RoundingMode defines how a value should be snapped
to the nearest preferred number.
*/
type RoundingMode int

const (
	// RoundingUnspecified is the zero value; treated as RoundNearest by all functions.
	RoundingUnspecified RoundingMode = iota

	// RoundNearest selects the closest preferred value.
	RoundNearest

	// RoundUp selects the next highest preferred value.
	RoundUp

	// RoundDown selects the next lowest preferred value.
	RoundDown
)

///////////////////////////////////////////////////////////////////////////////
// Base Decade Tables (Hybrid Approach)
///////////////////////////////////////////////////////////////////////////////

/*
eSeriesBase contains canonical IEC base decade values
for lower series.

These values exist in the normalized range:

	1.0 ≤ value < 10.0

Higher series (E48, E96, E192) are generated mathematically
to avoid maintaining very large constant tables.

This hybrid approach ensures:

  - Readability for small series
  - Maintainability for large series
  - Standards compliance
*/
var eSeriesBase = map[ESeries][]float64{
	E3: {
		1.0, 2.2, 4.7,
	},
	E6: {
		1.0, 1.5, 2.2, 3.3, 4.7, 6.8,
	},
	E12: {
		1.0, 1.2, 1.5, 1.8, 2.2, 2.7,
		3.3, 3.9, 4.7, 5.6, 6.8, 8.2,
	},
	E24: {
		1.0, 1.1, 1.2, 1.3, 1.5, 1.6,
		1.8, 2.0, 2.2, 2.4, 2.7, 3.0,
		3.3, 3.6, 3.9, 4.3, 4.7, 5.1,
		5.6, 6.2, 6.8, 7.5, 8.2, 9.1,
	},
}

func init() {
	eSeriesBase[E48] = generateESeries(48)
	eSeriesBase[E96] = generateESeries(96)
	eSeriesBase[E192] = generateESeries(192)
}

///////////////////////////////////////////////////////////////////////////////
// Series Generation
///////////////////////////////////////////////////////////////////////////////

/*
generateESeries produces normalized preferred values
for higher-resolution series (E48, E96, E192).

IEC 60063 defines preferred numbers as logarithmically spaced:

	value_n = 10^(n/N)

Where:

	N = number of values per decade
	n = 0 to N-1

These values are then rounded to a defined number of
significant digits to produce commercially usable numbers.

We round to 3 significant digits for higher series,
which matches common manufacturer rounding behavior.
*/
func generateESeries(N int) []float64 {
	values := make([]float64, 0, N)

	for n := range N {
		v := math.Pow(10, float64(n)/float64(N))
		v = roundToSignificant(v, 3)
		values = append(values, v)
	}

	return values
}

/*
roundToSignificant rounds a floating-point number
to the specified number of significant digits.

Example:

	roundToSignificant(1.23456, 3) → 1.23
	roundToSignificant(9.8765, 2) → 9.9

This ensures stable IEC-style rounding for generated series.
*/
func roundToSignificant(x float64, sig int) float64 {
	if x == 0 {
		return 0
	}

	scale := math.Pow(10, float64(sig)-math.Ceil(math.Log10(math.Abs(x))))
	return math.Round(x*scale) / scale
}

///////////////////////////////////////////////////////////////////////////////
// Internal Base Fetcher
///////////////////////////////////////////////////////////////////////////////

/*
baseValues returns the normalized base decade values
for the given E-series.

For small series (E3–E24), predefined canonical tables are used.
For larger series (E48–E192), values are generated mathematically.
*/
func baseValues(series ESeries) ([]float64, bool) {
	base, ok := eSeriesBase[series]
	return base, ok
}

///////////////////////////////////////////////////////////////////////////////
// Public API
///////////////////////////////////////////////////////////////////////////////

/*
NearestStandard snaps a resistance value to the nearest
preferred value in the specified E-series.

Process Overview:

 1. Normalize input value into its decade:
    value = normalized × 10^exponent

    Example:
    4700 → 4.7 × 10^3

 2. Compare normalized value to preferred numbers
    within 1.0–10.0 range.

3. Apply rounding rule (Nearest, Up, Down).

4. Re-scale back to original decade.

Example:

	NearestStandard(500, E24, RoundNearest)
	→ 510

	NearestStandard(500, E12, RoundNearest)
	→ 470

	NearestStandard(500, E12, RoundUp)
	→ 560

This function does NOT:
  - Validate manufacturability
  - Snap tolerance
  - Perform band encoding

It strictly performs IEC 60063 value selection.
*/
func NearestStandard(value float64, series ESeries, mode RoundingMode) (float64, error) {

	if value <= 0 {
		return 0, errors.New("value must be positive")
	}

	base, ok := baseValues(series)
	if !ok {
		return 0, errors.New("unsupported E-series")
	}

	// Determin decade component.
	exponent := math.Floor(math.Log10(value))

	// Normalize value into [1, 10) range.
	normalized := value / math.Pow(10, exponent)

	// Binary search: i is the first index where base[i] >= normalized.
	i := sort.SearchFloat64s(base, normalized)

	var best float64
	switch mode {
	case RoundUp:
		if i < len(base) {
			best = base[i]
		} else {
			result := base[0] * math.Pow(10, exponent+1)
			return roundToSignificant(result, 6), nil
		}
	case RoundDown:
		if i < len(base) && base[i] <= normalized {
			best = base[i]
		} else if i > 0 {
			best = base[i-1]
		} else {
			best = base[0]
		}
	default: // RoundNearest and RoundingUnspecified
		switch {
		case i == 0:
			best = base[0]
		case i == len(base):
			last := base[len(base)-1]
			next := base[0] * math.Pow(10, exponent+1)
			if math.Abs(normalized-next/math.Pow(10, exponent)) < math.Abs(normalized-last) {
				result := next
				return roundToSignificant(result, 6), nil
			}
			best = last
		default:
			lo, hi := base[i-1], base[i]
			if math.Abs(normalized-lo) <= math.Abs(normalized-hi) {
				best = lo
			} else {
				best = hi
			}
		}
	}

	result := best * math.Pow(10, exponent)
	return roundToSignificant(result, 6), nil
}
