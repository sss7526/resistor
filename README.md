# resistor

A Go library for working with fixed resistors. Converts between resistance values and visual encodings, snaps values to standard preferred series, and infers unknown properties from physical observations.

The library has no UI dependencies and is suitable for use in other Go projects, command-line tools, or WebAssembly modules.

---

## Contents

- [Library](#library)
- [CLI](#cli)
- [TUI](#tui)
- [Building](#building)
- [Testing](#testing)
- [License](#license)

---

## Library

```
go get github.com/sss7526/resistor
```

### Color Code Decoding and Encoding

Decode 4-band, 5-band, and 6-band resistors per IEC 60062:

```go
bands := []resistor.Color{resistor.Green, resistor.Brown, resistor.Brown, resistor.Gold}
spec, err := resistor.DecodeBands(bands)
// spec.ResistanceOhms == 510
// spec.TolerancePct  == 5
```

Encode a resistance and tolerance back to bands:

```go
spec := resistor.ResistorSpec{ResistanceOhms: 510, TolerancePct: 5}
bands, err := resistor.EncodeBands(spec)
// []Color{Green, Brown, Brown, Gold}
```

6-band resistors include a temperature coefficient in the sixth band:

```go
spec := resistor.ResistorSpec{
    ResistanceOhms: 4700,
    TolerancePct:   1,
    TempCoeffPPM:   50,
}
bands, err := resistor.EncodeBands(spec)
```

### SMD Markings

Decode surface-mount markings in 3-digit, 4-digit, R-notation, and EIA-96 formats:

```go
spec, err := resistor.DecodeSMD("472")   // 4700 ohms
spec, err := resistor.DecodeSMD("4R7")   // 4.7 ohms
spec, err := resistor.DecodeSMD("01C")   // EIA-96
```

Encode a resistance value to an SMD marking:

```go
marking, err := resistor.EncodeSMD(4700, resistor.SMDAuto)
// "472"
```

### E-Series Value Selection

Snap an arbitrary resistance to the nearest value in a standard IEC 60063 preferred series (E3 through E192):

```go
v, err := resistor.NearestStandard(487, resistor.E24, resistor.RoundNearest)
// 487 -> 510

v, err := resistor.NearestStandard(487, resistor.E12, resistor.RoundUp)
// 487 -> 560
```

### Standard Resistor Selection

Select a standard resistor for a target resistance. Returns the snapped value, color bands, and a record of any defaults that were applied:

```go
result, err := resistor.SelectStandardResistor(resistor.SelectionRequest{
    Resistance: 487,
    Series:     resistor.E24,
})
// result.SelectedResistance == 510
// result.Bands              == [Green, Brown, Brown, Gold]
// result.Assumptions        == ["Tolerance defaulted to +/-5%", ...]
```

Unspecified fields default to E24, 5% tolerance, and nearest rounding. All defaults are recorded in `result.Assumptions`.

### Physical Inference

Estimate unknown properties from a physical observation. Accepts any combination of color bands, SMD marking, body color, length, and package type:

```go
result, err := resistor.InferResistor(resistor.ObservedResistor{
    Bands:     []resistor.Color{resistor.Brown, resistor.Black, resistor.Red, resistor.Gold},
    BodyColor: resistor.Blue,
    LengthMM:  6.3,
})
// result.Spec.ResistanceOhms == 1000
// result.Spec.PowerWatts     == 0.25  (inferred from length)
// result.Spec.Type           == "metal_film"  (inferred from body color)
// result.Meta.Confidence     == 0.92
// result.Meta.Assumptions    == ["Blue body assumed metal film", "Length 5-7mm assumed 1/4W"]
```

Deterministic facts (decoded from bands or markings) always take precedence over heuristic estimates. Confidence is a value in [0.0, 1.0] computed as a weighted average of the rules that fired.

### Engineering Analysis

Analyze a resistor under specified electrical conditions:

```go
report, err := resistor.AnalyzeResistor(resistor.AnalysisInput{
    Spec: resistor.ResistorSpec{
        ResistanceOhms: 100,
        PowerWatts:     0.25,
        TolerancePct:   5,
    },
    AppliedVoltage: 10,
})
// report.Current                == 0.1 A
// report.PowerDissipation       == 0.1 W
// report.DeratedSafePower       == 0.125 W  (50% derating)
// report.WorstCaseResistanceMin == 95 ohms
// report.WorstCaseResistanceMax == 105 ohms
// report.Warnings               == [{Caution, "Power dissipation exceeds recommended 50% derated threshold"}]
```

Either `AppliedVoltage` or `AppliedCurrent` may be provided. If both are given, consistency with Ohm's Law is checked. Missing inputs produce structured warnings rather than errors.

---

## CLI

The CLI provides non-interactive access to all library operations.

### Installation

```
go install github.com/sss7526/resistor/cmd/resistor-cli@latest
```

Or build from source:

```
make build-cli
./bin/resistor-cli --help
```

### Commands

**Select a standard resistor value:**

```
resistor-cli select 487
resistor-cli select 487 --series E12
resistor-cli select 487 --tol 1 --round up
```

**Infer properties from physical observations:**

```
resistor-cli infer --bands brown,black,red,gold
resistor-cli infer --bands brown,black,red,gold --body blue --length 6.3
resistor-cli infer --smd 472 --pkg 0603
```

**Analyze under electrical conditions:**

```
resistor-cli analyze --r 100 --v 10 --pwr 0.25 --tol 5
resistor-cli analyze --r 100 --i 0.1
```

**Decode or encode SMD markings:**

```
resistor-cli smd decode 472
resistor-cli smd decode 4R7
resistor-cli smd encode 4700
```

All commands accept `--json` for machine-readable output.

---

## TUI

The TUI provides an interactive terminal interface for the same operations.

```
make build-tui
./bin/resistor-tui
```

Or install:

```
go install github.com/sss7526/resistor/cmd/resistor-tui@latest
```

Navigate with arrow keys, confirm with Enter, and return to the menu with Escape. Press `q` or `Ctrl+C` to quit.

---

## Building

Requires Go 1.21 or later.

```
make build        # build both binaries to bin/
make build-cli    # build CLI only
make build-tui    # build TUI only
```

---

## Testing

```
make test         # unit tests
make test-cli     # CLI integration tests
make test-all     # both
make smoke        # end-to-end smoke tests against the built CLI binary
```

Fuzz testing:

```
make fuzz         # run all fuzz targets for 10 seconds each
make fuzz FUZZTIME=60s
```

---

## License

MIT. See [LICENSE](LICENSE).
