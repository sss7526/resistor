package resistor

import (
	"errors"
	"math"
)

///////////////////////////////////////////////////////////////////////////////
// Errors
///////////////////////////////////////////////////////////////////////////////

// ErrInvalidBandCount indicates that the provided band slice
// does not contain 4, 5, or 6 entries.
//
// 4-band format:
//
//	digit, digit, multiplier, tolerance
//
// 5-band format:
//
//	digit, digit, digit, multiplier, tolerance
//
// 6-band format:
//
//	digit, digit, digit, multiplier, tolerance, temp coefficient
var ErrInvalidBandCount = errors.New("invalid number of bands (must be 4, 5, or 6)")

// ErrInvalidDigitColor indicates that a band expected to represent
// a significant digit does not map to a valid digit color (0–9).
var ErrInvalidDigitColor = errors.New("invalid digit color")

// ErrInvalidMultiplier indicates the multiplier band color
// does not map to a valid multiplier value.
var ErrInvalidMultiplier = errors.New("invalid multiplier color")

// ErrInvalidTolerance indicates the tolerance band color
// does not map to a known tolerance value.
var ErrInvalidTolerance = errors.New("invalid tolerance color")

// ErrUnencodableValue indicates that the provided resistance
// cannot be represented exactly in standard 4-band or 5-band format.
//
// This happens when:
//   - The resistance requires more than 2 or 3 significant digits
//   - The resistance does not divide cleanly by any defined multiplier
//   - The tolerance cannot be represented by a standard color
var ErrUnencodableValue = errors.New("resistance cannot be encoded in 4 or 5 band format")

///////////////////////////////////////////////////////////////////////////////
// Decoding
///////////////////////////////////////////////////////////////////////////////

/*
DecodeBands converts a slice of 4, 5, or 6 color bands into a ResistorSpec.

Only the following fields of ResistorSpec are populated:

  - ResistanceOhms
  - TolerancePct
  - TempCoeffPPM (6-band only)

No inference is performed.
No power rating is assumed.

Resistor Band Structure (IEC 60062):

4-band resistor:

	Band 1: First significant digit
	Band 2: Second significant digit
	Band 3: Multiplier (power of ten scaling factor)
	Band 4: Tolerance

5-band resistor:

	Band 1: First significant digit
	Band 2: Second significant digit
	Band 3: Third significant digit
	Band 4: Multiplier
	Band 5: Tolerance

6-band resistor:

	Band 1: First significant digit
	Band 2: Second significant digit
	Band 3: Third significant digit
	Band 4: Multiplier
	Band 5: Tolerance
	Band 6: Temperature coefficient

Example (4-band):

	Green, Brown, Brown, Gold

	Digits: 5 1
	Multiplier: 10
	→ 51 × 10 = 510Ω
	Tolerance: ±5%

Mathematical Model:

Resistance = (significant digits combined as integer) × multiplier

Example:

	Digits: [5, 1]
	Combined: 51
	Multiplier: 10
	Result: 51 × 10 = 510Ω

This function strictly follows that model.
*/
func DecodeBands(bands []Color) (ResistorSpec, error) {
	var spec ResistorSpec

	if len(bands) != 4 && len(bands) != 5 && len(bands) != 6 {
		return spec, ErrInvalidBandCount
	}

	var digits []int
	var multiplier float64
	var tolerance float64
	var ok bool

	if len(bands) == 4 {
		// 4-band: 2 digits + multiplier + tolerance

		d1, ok1 := DigitValue[bands[0]]
		d2, ok2 := DigitValue[bands[1]]
		if !ok1 || !ok2 {
			return spec, ErrInvalidDigitColor
		}
		digits = []int{d1, d2}

		multiplier, ok = MultiplierValue[bands[2]]
		if !ok {
			return spec, ErrInvalidMultiplier
		}

		tolerance, ok = ToleranceValue[bands[3]]
		if !ok {
			return spec, ErrInvalidTolerance
		}

	} else {
		// 5-band or 6 band: 3 digits + multiplier + tolerance

		d1, ok1 := DigitValue[bands[0]]
		d2, ok2 := DigitValue[bands[1]]
		d3, ok3 := DigitValue[bands[2]]
		if !ok1 || !ok2 || !ok3 {
			return spec, ErrInvalidDigitColor
		}
		digits = []int{d1, d2, d3}

		multiplier, ok = MultiplierValue[bands[3]]
		if !ok {
			return spec, ErrInvalidMultiplier
		}

		tolerance, ok = ToleranceValue[bands[4]]
		if !ok {
			return spec, ErrInvalidTolerance
		}
	}

	// Combine digits into a base integer.
	// Example: [5,1,2] becomes 512.
	value := 0
	for _, d := range digits {
		value = value*10 + d
	}

	// Apply multiplier to compute final resistance.
	spec.ResistanceOhms = float64(value) * multiplier
	spec.TolerancePct = tolerance

	// If 6-band, decode temperature coefficient
	if len(bands) == 6 {
		if ppm, ok := TempCoeffValue[bands[5]]; ok {
			spec.TempCoeffPPM = ppm
		}
	}

	return spec, nil
}

///////////////////////////////////////////////////////////////////////////////
// Encoding
///////////////////////////////////////////////////////////////////////////////

/*
EncodeBands converts a ResistorSpec into IEC color bands.

Behavior:

4-band:

	Used when tolerance > 2% and TempCoeffPPM == 0

5-band:

	Used when tolerance ≤ 2% and TempCoeffPPM == 0

6-band:

	Used when TempCoeffPPM != 0

TempCoeffPPM must match a defined TempCoeffValue entry exactly.
*/
func EncodeBands(spec ResistorSpec) ([]Color, error) {

	if spec.ResistanceOhms <= 0 {
		return nil, ErrUnencodableValue
	}

	// 6-band encoding if tempco provided
	if spec.TempCoeffPPM != 0 {
		return encodeSixBand(spec)
	}

	// 5-band for precision
	if spec.TolerancePct <= 2.0 {
		return encodeFiveBand(spec.ResistanceOhms, spec.TolerancePct)
	}

	// Default 4-band
	return encodeFourBand(spec.ResistanceOhms, spec.TolerancePct)
}

/*
EncodeBandsSimple preserves backward compatibility.

Encodes resistance and tolerance without tempco.
*/
func EncodeBandsSimple(resistance, tolerance float64) ([]Color, error) {
	return EncodeBands(ResistorSpec{
		ResistanceOhms: resistance,
		TolerancePct:   tolerance,
	})
}

///////////////////////////////////////////////////////////////////////////////
// 4-Band Encoding
///////////////////////////////////////////////////////////////////////////////

func encodeFourBand(resistance float64, tolerance float64) ([]Color, error) {

	for multiplierColor, multiplier := range MultiplierValue {

		// Reverse the multiplier operation:
		// resistance = digits × multiplier
		// digits = resistance / multiplier
		value := resistance / multiplier

		// 4-band requires exactly 2 digits:
		// valid range: 10–99
		if value < 10 || value >= 100 {
			continue
		}

		// Must be whole number (no fractional digits allowed)
		if math.Mod(value, 1) != 0 {
			continue
		}

		intVal := int(value)
		d1 := intVal / 10
		d2 := intVal % 10

		c1, ok1 := DigitColor[d1]
		c2, ok2 := DigitColor[d2]
		tolColor, okTol := ToleranceColor[tolerance]

		if !ok1 || !ok2 || !okTol {
			return nil, ErrUnencodableValue
		}

		return []Color{
			c1,
			c2,
			multiplierColor,
			tolColor,
		}, nil
	}

	return nil, ErrUnencodableValue
}

///////////////////////////////////////////////////////////////////////////////
// 5-Band Encoding
///////////////////////////////////////////////////////////////////////////////

func encodeFiveBand(resistance float64, tolerance float64) ([]Color, error) {

	for multiplierColor, multiplier := range MultiplierValue {

		value := resistance / multiplier

		// 5-band requires exactly 3 digits:
		// valid range: 100–999
		if value < 100 || value >= 1000 {
			continue
		}

		if math.Mod(value, 1) != 0 {
			continue
		}

		intVal := int(value)
		d1 := intVal / 100
		d2 := (intVal / 10) % 10
		d3 := intVal % 10

		c1, ok1 := DigitColor[d1]
		c2, ok2 := DigitColor[d2]
		c3, ok3 := DigitColor[d3]
		tolColor, okTol := ToleranceColor[tolerance]

		if !ok1 || !ok2 || !ok3 || !okTol {
			return nil, ErrUnencodableValue
		}

		return []Color{
			c1,
			c2,
			c3,
			multiplierColor,
			tolColor,
		}, nil
	}

	return nil, ErrUnencodableValue
}

func encodeSixBand(spec ResistorSpec) ([]Color, error) {

	bands, err := encodeFiveBand(spec.ResistanceOhms, spec.TolerancePct)
	if err != nil {
		return nil, err
	}

	tempColor, ok := tempCoeffToColor(spec.TempCoeffPPM)
	if !ok {
		return nil, ErrUnencodableValue
	}

	return append(bands, tempColor), nil
}

func tempCoeffToColor(ppm int) (Color, bool) {
	c, ok := TempCoeffColor[ppm]
	return c, ok
}
