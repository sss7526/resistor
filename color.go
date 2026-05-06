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

	None Color = "none" // used when no tolerance band exists (±20%)
)
