package resistor

// DigitValue maps band color → digit.
var DigitValue = map[Color]int{
	Black:  0,
	Brown:  1,
	Red:    2,
	Orange: 3,
	Yellow: 4,
	Green:  5,
	Blue:   6,
	Violet: 7,
	Grey:   8,
	White:  9,
}

// DigitColor maps digit → band color.
var DigitColor = map[int]Color{
	0: Black,
	1: Brown,
	2: Red,
	3: Orange,
	4: Yellow,
	5: Green,
	6: Blue,
	7: Violet,
	8: Grey,
	9: White,
}

// MultiplierValue maps band color → multiplier factor.
var MultiplierValue = map[Color]float64{
	Black:  1,
	Brown:  10,
	Red:    100,
	Orange: 1_000,
	Yellow: 10_000,
	Green:  100_000,
	Blue:   1_000_000,
	Violet: 10_000_000,
	Grey:   100_000_000,
	White:  1_000_000_000,
	Gold:   0.1,
	Silver: 0.01,
}

// MultiplierColor maps multiplier factor → band color.
var MultiplierColor = map[float64]Color{
	1:             Black,
	10:            Brown,
	100:           Red,
	1_000:         Orange,
	10_000:        Yellow,
	100_000:       Green,
	1_000_000:     Blue,
	10_000_000:    Violet,
	100_000_000:   Grey,
	1_000_000_000: White,
	0.1:           Gold,
	0.01:          Silver,
}

// ToleranceValue maps band color → tolerance percentage.
var ToleranceValue = map[Color]float64{
	Brown:  1.0,
	Red:    2.0,
	Green:  0.5,
	Blue:   0.25,
	Violet: 0.1,
	Grey:   0.05,
	Gold:   5.0,
	Silver: 10.0,
	None:   20.0,
}

// ToleranceColor maps tolerance percentage → band color.
var ToleranceColor = map[float64]Color{
	1.0:  Brown,
	2.0:  Red,
	0.5:  Green,
	0.25: Blue,
	0.1:  Violet,
	0.05: Grey,
	5.0:  Gold,
	10.0: Silver,
	20.0: None,
}

// TempCoeffValue maps 6th band color → temperature coefficient in ppm/°C.
// Defined according to IEC 60062.
var TempCoeffValue = map[Color]int{
	Brown:  100,
	Red:    50,
	Orange: 15,
	Yellow: 25,
	Blue:   10,
	Violet: 5,
}

// TempCoeffColor maps ppm value → color.
var TempCoeffColor = map[int]Color{
	100: Brown,
	50:  Red,
	15:  Orange,
	25:  Yellow,
	10:  Blue,
	5:   Violet,
}
