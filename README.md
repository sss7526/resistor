# resistor

A Go library for working with fixed resistors. Decode and encode color band and SMD markings, snap values to IEC preferred series, and infer unknown properties from physical observations.

No UI dependencies. Import it as a library, or use the included CLI, TUI, and web app.

## Contents

- [Library](#library)
- [CLI](#cli)
- [TUI](#tui)
- [Web](#web)
- [Deployment](#deployment)
- [Development](#development)
- [License](#license)

---

## Library

```
go get github.com/sss7526/resistor
```

### Color band decode and encode

Supports 4-band, 5-band, and 6-band resistors per IEC 60062:

```go
bands := []resistor.Color{resistor.Green, resistor.Brown, resistor.Brown, resistor.Gold}
spec, err := resistor.DecodeBands(bands)
// spec.ResistanceOhms == 510, spec.TolerancePct == 5

spec := resistor.ResistorSpec{ResistanceOhms: 510, TolerancePct: 5}
bands, err := resistor.EncodeBands(spec)
// []Color{Green, Brown, Brown, Gold}
```

6-band resistors include a temperature coefficient in the sixth band.

### SMD markings

Supports 3-digit, 4-digit, R-notation, and EIA-96 formats:

```go
spec, err := resistor.DecodeSMD("472")   // 4700 ohms
spec, err := resistor.DecodeSMD("4R7")   // 4.7 ohms
spec, err := resistor.DecodeSMD("01C")   // EIA-96

marking, err := resistor.EncodeSMD(4700, resistor.SMDAuto)
// "472"
```

### E-series value selection

Snap an arbitrary resistance to the nearest IEC 60063 preferred series value (E3 through E192):

```go
v, err := resistor.NearestStandard(487, resistor.E24, resistor.RoundNearest) // 510
v, err := resistor.NearestStandard(487, resistor.E12, resistor.RoundUp)      // 560
```

### Standard resistor selection

Returns the snapped value, color bands, and all defaults applied:

```go
result, err := resistor.SelectStandardResistor(resistor.SelectionRequest{
    Resistance: 487,
    Series:     resistor.E24,
})
// result.SelectedResistance == 510
// result.Bands              == [Green, Brown, Brown, Gold]
// result.Assumptions        == ["Tolerance defaulted to +/-5%", ...]
```

Unspecified fields default to E24, 5% tolerance, and nearest rounding.

### Physical inference

Estimate unknown properties from any combination of bands, SMD marking, body color, length, and package type:

```go
result, err := resistor.InferResistor(resistor.ObservedResistor{
    Bands:     []resistor.Color{resistor.Brown, resistor.Black, resistor.Red, resistor.Gold},
    BodyColor: resistor.Blue,
    LengthMM:  6.3,
})
// result.Spec.ResistanceOhms == 1000
// result.Spec.PowerWatts     == 0.25         (inferred from length)
// result.Spec.Type           == "metal_film" (inferred from body color)
// result.Meta.Confidence     == 0.92
// result.Meta.Assumptions    == ["Blue body assumed metal film", ...]
```

Decoded facts always take precedence over heuristic estimates. Confidence is a value in [0.0, 1.0].

### Engineering analysis

Accepts voltage or current (or both, checked for Ohm's Law consistency):

```go
report, err := resistor.AnalyzeResistor(resistor.AnalysisInput{
    Spec:           resistor.ResistorSpec{ResistanceOhms: 100, PowerWatts: 0.25, TolerancePct: 5},
    AppliedVoltage: 10,
})
// report.Current                == 0.1 A
// report.PowerDissipation       == 0.1 W
// report.DeratedSafePower       == 0.125 W
// report.WorstCaseResistanceMin == 95 ohms
// report.WorstCaseResistanceMax == 105 ohms
// report.Warnings               == [{Caution, "Power dissipation exceeds 50% derated threshold"}]
```

---

## CLI

```
go install github.com/sss7526/resistor/cmd/resistor-cli@latest
```

All commands accept `--json` for machine-readable output.

**Select a standard value:**
```
resistor-cli select 487
resistor-cli select 487 --series E12 --tol 1 --round up
```

**Infer from physical observations:**
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

**Decode and encode SMD markings:**
```
resistor-cli smd decode 472
resistor-cli smd decode 4R7
resistor-cli smd encode 4700
```

---

## TUI

An interactive terminal interface for all operations.

```
go install github.com/sss7526/resistor/cmd/resistor-tui@latest
```

Navigate with arrow keys, confirm with Enter, return to menu with Escape, quit with `q` or `Ctrl+C`.

---

## Web

The WASM module exposes all library operations as browser-callable JavaScript functions. The included server embeds the compiled WASM and all static assets into a single binary.

### Build

```
make build-wasm           # standard Go WASM (~3.4 MB uncompressed)
make build-server         # server + standard WASM
make build-server-tinygo  # server + TinyGo WASM (~430 KB gzip, ~1.1 MB uncompressed)
```

TinyGo must be installed for the TinyGo targets. See [tinygo.org/getting-started/install](https://tinygo.org/getting-started/install).

Install to `GOPATH/bin`:
```
make install-server
make install-server-tinygo
```

### Run

```
./bin/resistor-server               # listens on :8080
./bin/resistor-server --addr :9000  # custom port
RESISTOR_ADDR=:9000 ./bin/resistor-server
```

The server sets strict security headers (CSP with per-request nonce, COEP, COOP, X-Frame-Options) and is designed to sit behind a TLS-terminating reverse proxy.

### JavaScript API

All functions live on the `resistor` global and return `{ok: true, value: ...}` on success or `{ok: false, error: "..."}` on failure. Inputs are JSON strings.

| Function | Input | Returns |
|---|---|---|
| `decodeBands` | `'["green","brown","brown","gold"]'` | `ResistorSpec` |
| `encodeBands` | `'{"resistanceOhms":510,"tolerancePct":5}'` | `Color[]` |
| `decodeSMD` | `"472"` | `ResistorSpec` |
| `encodeSMD` | `'{"resistance":4700,"mode":"auto"}'` | marking string |
| `nearestStandard` | `'{"value":487,"series":"E24","mode":"nearest"}'` | number |
| `selectStandardResistor` | `'{"resistance":487}'` | `SelectionResult` |
| `inferResistor` | `'{"bands":[...],"bodyColor":"blue","lengthMM":6.3}'` | `InferenceResult` |
| `analyzeResistor` | `'{"spec":{...},"appliedVoltage":10}'` | `AnalysisReport` |

---

## Deployment

The server binary is self-contained and designed to run behind Caddy (or any reverse proxy). Caddy handles TLS automatically via Let's Encrypt.

### Docker

```
# Standard Go WASM
docker build -t resistor .
docker run --rm -p 8080:8080 resistor

# TinyGo WASM (~430 KB gzip)
docker build --build-arg WASM=tinygo -t resistor .
```

The final image is built on `scratch` and runs as UID 10001.

### Docker Compose + Caddy

Set `RESISTOR_HOST` to your public domain and start:

```
RESISTOR_HOST=resistor.example.com docker compose up -d
```

Caddy provisions TLS automatically. For local testing (HTTP only, no cert needed):

```
docker compose up
```

To use TinyGo WASM:

```
WASM=tinygo RESISTOR_HOST=resistor.example.com docker compose up -d
```

Certificate data is persisted in the `caddy_data` named volume. Edit `Caddyfile` before starting Compose to add rate limiting, authentication, or other directives.

---

## Development

Requires Go 1.26 or later.

**Build:**
```
make build        # CLI + TUI to bin/
make build-cli
make build-tui
make build-server
```

**Test:**
```
make test         # unit tests
make test-cli     # CLI integration tests
make test-all     # both
make smoke        # end-to-end smoke tests against the built binary
```

**Fuzz:**
```
make fuzz
make fuzz FUZZTIME=60s
```

**Benchmark:**
```
go test -bench=. -benchtime=1s ./...
go test -bench=BenchmarkInferResistor -benchmem ./...
```

**Lint and vulnerability check:**
```
make lint
make vuln
```

---

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

---

## License

MIT. See [LICENSE](LICENSE).
