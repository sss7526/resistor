package resistor

// DigitValue maps color to significant digit (0-9)
var DigitValue = map[Color]int{
	Black: 0,
	Brown: 1,
	Red: 2,
	Orange: 3,
	Yellow: 4,
	Green: 5,
	Blue: 6,
	Violet: 7,
	Grey: 8,
	White: 9,
}

// MultiplierValue maps color to multiplier value.
var MultiplierValue = map[Color]float64{
	Black: 1,
	Brown: 10,
	Red: 100,
	Orange: 1_000,
	Yellow: 10_000,
	Green: 100_000,
	Blue: 1_000_000,
	Violet: 10_000_000,
	Grey: 100_000_000,
	White: 1_000_000_000,

	Gold: 0.1,
	Silver: 0.01,
}

// ToleranceValue maps color to tolerance percentage.
var ToleranceValue = map[Color]float64{
	Brown: 1.0,
	Red: 2.0,
	Green: 0.5,
	Blue: 0.25,
	Violet: 0.1,
	Grey: 0.05,
	Gold: 5.0,
	Silver: 10.0,
	None: 20.0,
}

// TempCoeffValue maps color to temperature coefficient in ppm/°C.
var TempCoeffPPM = map[Color]int{
	Brown: 100,
	Red: 50,
	Orange: 15,
	Yellow: 25,
	Blue: 10,
	Violet: 5,
}