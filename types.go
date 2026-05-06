package resistor

// PackageType represents physical mounting style.
type PackageType string

const (
	ThroughHole PackageType = "through_hole"
	SMD         PackageType = "smd"
	Axial       PackageType = "axial"
	Radial      PackageType = "radial"
	UnknownPKG  PackageType = "unkown"
)

// ResistorType represents construction material
type ResistorType string

const (
	CarbonFilm ResistorType = "carbon_film"
	MetalFilm  ResistorType = "metal_film"
	ThickFilm  ResistorType = "thick_film"
	ThinFilm   ResistorType = "thin_film"
	Wirewound  ResistorType = "wirewound"
	UnkownType ResistorType = "unknown"
)

// ResistorSpec represent a fully specified resistor.
// This is for known, deterministic data.
type ResistorSpec struct {
	ResistanceOhms float64
	TolerancePct   float64
	PowerWatts     float64
	TempCoeffPPM   int
	Package        PackageType
	Type           ResistorType
}

// VisualProfile represents a visual encoding of a resistor
type VisualProfile struct {
	Bands     []Color
	BodyColor Color
	SMDCode   string
}

// ObservedResistor represents incomplete data collected
// from a physical resistor without packaging.
type ObservedResistor struct {
	Bands     []Color
	BodyColor Color
	LengthMM  float64
	Package   PackageType
	Marking   string // e.g. SMD code "472"
}

// InferenceMeta represents metadat for future inference stages.
// Included now so it is part of stable domain modeling.
type InferenceMeta struct {
	Assumptions []string
	Confidence  float64
}
