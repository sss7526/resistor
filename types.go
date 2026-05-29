package resistor

// PackageType represents the physical mounting style of a resistor.
type PackageType string

const (
	// ThroughHole is a standard leaded resistor for PCB through-hole mounting.
	ThroughHole PackageType = "through_hole"

	// SMD is a generic surface-mount resistor with unspecified case size.
	SMD PackageType = "smd"

	// SMD0402 is a 1.0×0.5 mm surface-mount case (1/16 W typical).
	SMD0402 PackageType = "smd_0402"

	// SMD0603 is a 1.6×0.8 mm surface-mount case (0.1 W typical).
	SMD0603 PackageType = "smd_0603"

	// SMD0805 is a 2.0×1.25 mm surface-mount case (0.125 W typical).
	SMD0805 PackageType = "smd_0805"

	// SMD1206 is a 3.2×1.6 mm surface-mount case (0.25 W typical).
	SMD1206 PackageType = "smd_1206"

	// Axial is a leaded resistor with leads extending from both ends along the axis.
	Axial PackageType = "axial"

	// Radial is a leaded resistor with leads extending from the same end.
	Radial PackageType = "radial"

	// UnknownPKG is used when the package type cannot be determined.
	UnknownPKG PackageType = "unknown"
)

// ResistorType represents the construction technology of a resistor.
type ResistorType string

const (
	// CarbonFilm is a carbon film resistor; common, low cost, moderate tolerance.
	CarbonFilm ResistorType = "carbon_film"

	// MetalFilm is a metal film resistor; tight tolerance, low noise, blue body typical.
	MetalFilm ResistorType = "metal_film"

	// ThickFilm is a thick film SMD resistor; standard surface-mount construction.
	ThickFilm ResistorType = "thick_film"

	// ThinFilm is a thin film SMD resistor; high precision, low temperature coefficient.
	ThinFilm ResistorType = "thin_film"

	// Wirewound is a wire-wound resistor; high power, high precision, inductive.
	Wirewound ResistorType = "wirewound"

	// MetalOxide is a metal oxide resistor; high stability, flame-resistant, green body typical.
	MetalOxide ResistorType = "metal_oxide"

	// UnkownType is used when the resistor construction type cannot be determined.
	UnkownType ResistorType = "unknown"
)

// ResistorSpec represents a fully specified resistor with known, deterministic values.
type ResistorSpec struct {
	// ResistanceOhms is the nominal resistance in ohms.
	ResistanceOhms float64

	// TolerancePct is the tolerance as a percentage (e.g. 5 for ±5%).
	TolerancePct float64

	// PowerWatts is the rated power dissipation in watts.
	PowerWatts float64

	// TempCoeffPPM is the temperature coefficient in ppm/°C (6-band resistors only).
	TempCoeffPPM int

	// Package is the physical mounting style.
	Package PackageType

	// Type is the construction technology.
	Type ResistorType
}

// VisualProfile represents a visual encoding of a resistor's markings.
type VisualProfile struct {
	// Bands holds the ordered color bands (4, 5, or 6 entries).
	Bands []Color

	// BodyColor is the primary body color of the resistor package.
	BodyColor Color

	// SMDCode is the printed surface-mount marking (e.g. "472", "01C").
	SMDCode string
}

// ObservedResistor represents incomplete data collected from a physical
// resistor without reference to its original packaging or datasheet.
// All fields are optional; provide whichever are observable.
type ObservedResistor struct {
	// Bands holds observed color bands. Valid counts are 4, 5, or 6.
	Bands []Color

	// BodyColor is the body color of the resistor (used for type inference).
	BodyColor Color

	// LengthMM is the physical body length in millimeters (used for power inference).
	LengthMM float64

	// Package is the observed package type (used for power and voltage inference).
	Package PackageType

	// Marking is the printed SMD code, if visible (e.g. "472", "4R7", "01C").
	Marking string
}

// InferenceMeta contains metadata produced by the inference engine.
type InferenceMeta struct {
	// Assumptions lists the human-readable reasoning steps applied during inference.
	Assumptions []string

	// Confidence is the aggregated confidence score in [0.0, 1.0].
	// Higher values indicate more corroborating evidence.
	Confidence float64
}
