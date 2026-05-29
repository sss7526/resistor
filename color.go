package resistor

// Color represents a standard IEC resistor band color.
type Color string

const (
	Black  Color = "black"
	Brown  Color = "brown"
	Red    Color = "red"
	Orange Color = "orange"
	Yellow Color = "yellow"
	Green  Color = "green"
	Blue   Color = "blue"
	Violet Color = "violet"
	Grey   Color = "grey"
	White  Color = "white"

	Gold   Color = "gold"
	Silver Color = "silver"

	Beige Color = "beige"
	Tan   Color = "tan"

	None Color = "none" // used when no tolerance band exists (±20%)
)

// String returns the color name as a plain string.
func (c Color) String() string {
	return string(c)
}

/*
BodyColors returns colors meaningful for resistor body inference.

These are distinct from band digit colors.
*/
func BodyColors() []Color {
	return []Color{
		Blue,
		Green,
		Beige,
		Tan,
	}
}
