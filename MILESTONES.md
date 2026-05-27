# Project Vision

> A Go library that converts between resistor specifications and visual encodings, and performs structured, assumption-aware inference for unknown resistors.

Core principles:
- Deterministic where standards exist
- Explicit assumptions where inference is required
- No UI dependencies in the core library
- Safe defaults for hobbyists

Delivery targets:
- Core library: importable Go package for use in other projects
- CLI: non-interactive command-line toolkit
- TUI: interactive terminal application
- WASM: client-side web toolkit; server responsibility is serving static files only

---

# Milestone 0 — Domain & Standards Foundation

**Goal:** Establish the complete domain model before implementing any logic.

### Deliverables:
- Enumerations for:
  - Color
  - Tolerance
  - Multiplier
  - Tempco
  - Package type
  - Resistor type
- Standard maps:
  - IEC 60062 color mapping
  - Tolerance mapping
  - Multiplier mapping
- Core data structures:
  - `ResistorSpec`
  - `VisualProfile`
  - `ObservedResistor`
  - `InferenceMeta`

### Done When:
- Any standard through-hole resistor can be represented in structured form.
- No inference logic implemented.
- No E-series logic implemented.


---

# Milestone 1 — Deterministic Color Code Engine

**Goal:** Fully reversible 4-band and 5-band resistor support.

### Deliverables:

#### 1 Decode
Input: color bands  
Output: resistance + tolerance

#### 2 Encode
Input: resistance + tolerance  
Output: correct band colors

#### 3 Validation
- Error on invalid band count
- Error on impossible combinations

### Out of Scope:
- Inference
- Power rating guessing
- Body color interpretation
- E-series snapping

### Done When:
- `DecodeBands()` and `EncodeBands()` are fully unit tested.
- Pure deterministic transformations only.


---

# Milestone 2 — E-Series Engine (IEC 60063)

**Goal:** Snap arbitrary resistance values to the nearest standard preferred value.

### Deliverables:

- E6
- E12
- E24 (minimum required)
- Configurable rounding:
  - Nearest
  - Round up
  - Round down

### API:

```go
NearestStandard(value float64, series ESeries, mode RoundingMode) float64
```

### Done When:
- 500Ω snaps to 510Ω in E24
- 480Ω snaps correctly
- Decades scale correctly

---

# Milestone 3 — Spec → Visual Engine (Production-Useful State)

**Goal:** Full forward transformation from engineering intent to color bands.

### Input:
- Desired resistance
- Optional tolerance
- Optional series preference

### Output:
- Snapped standard value
- Band representation
- Structured assumptions

### Example Output:

```json
{
  "input": 500,
  "snapped_to": 510,
  "series": "E24",
  "bands": ["green","brown","brown","gold"],
  "assumptions": ["Tolerance defaulted to ±5%"]
}
```

### Done When:
- Correct band colors generated for any valid input.
- All defaults and assumptions are explicitly recorded in the result.

---

# Milestone 4 — SMD Decoder

**Goal:** Surface-mount resistor marking support.

### Deliverables:
- 3-digit decoding
- 4-digit decoding
- R notation decoding (4R7)
- EIA-96 decoding (optional but high value)

### Done When:
- Input "472" → 4700Ω
- Input "01C" (EIA-96) decodes correctly


---

# Milestone 5 — Physical Inference Engine (Assumption-Based)

**Goal:** Estimate unknown resistor properties from physical observations.

### Inference Inputs:
- Body color
- Measured length (mm)
- Band count
- Tolerance band
- Package guess

### Inference Outputs:
- Estimated power rating
- Estimated resistor type
- Estimated voltage rating (conservative)
- Confidence score
- Assumptions list

### Design Constraint:
All inference must:
- Declare assumptions
- Provide confidence score
- Never overwrite deterministic data

### Example:

```json
{
  "resistance": 1000,
  "tolerance": 5,
  "power_estimate": 0.25,
  "type_estimate": "metal_film",
  "confidence": 0.78,
  "assumptions": [
    "Blue body assumed metal film",
    "6.3mm length assumed 1/4W"
  ]
}
```


---

# Milestone 6 — Safety & Engineering Enhancements

**Goal:** Electrical analysis for a resistor under specified operating conditions.

### Features:
- Power dissipation calculator
- Derating recommendation (50% rule)
- Voltage drop estimator
- Tolerance worst-case bounds
- Warning system

---

# Milestone 7 — Confidence & Scoring Model Refinement

**Goal:** Replace ad hoc confidence rules with a formal, weighted rule type.

```go
type InferenceRule struct {
    Weight    float64
    Condition func(...)
}
```

---

# Milestone 8 — Core Library Stability

**Goal:** API stability and test coverage sufficient for a stable importable release.

### Deliverables:
- 100% deterministic coverage tests
- Fuzz testing on all decode paths
- Clear README documenting assumptions, API surface, and usage examples
- Stable public API with versioned module
- Benchmarks on hot paths (NearestStandard, DecodeBands, InferResistor)

### Done When:
- ✓ No exported symbols are expected to change.
- ✓ All deterministic operations have exhaustive unit tests.
- ✓ Public API is documented well enough for a third-party consumer.
- ✓ Benchmarks on hot paths (`NearestStandard`, `DecodeBands`, `InferResistor`).

---

# Milestone 9 — CLI ✓

**Goal:** A complete, tested command-line interface covering all core library operations.

### Commands:

| Command | Description |
|---|---|
| `select [resistance]` | Snap to nearest standard value and encode color bands |
| `infer` | Infer properties from bands, SMD marking, or physical observations |
| `analyze` | Electrical analysis under specified voltage or current conditions |
| `smd decode [marking]` | Decode SMD marking to resistance |
| `smd encode [resistance]` | Encode resistance to SMD marking |
| `version` | Print version |

All commands support `--json` for machine-readable output.

### Done When:
- Human-readable output for all commands reflects the full data returned by the library,
  with no gap between what is computed and what is displayed.
- Integration tests cover all commands, all flags, and error paths.
- `--json` output is consistent and machine-parseable across all commands.

---

# Milestone 10 — TUI

**Goal:** An interactive terminal application covering all four toolkit operations.

### Views:

| View | Status |
|---|---|
| Menu | Complete |
| Select Resistor | Complete |
| Infer Resistor | Complete |
| Analyze Resistor | Complete |
| SMD Tools | Complete |

### Resolved Stability Issues:
- `huh` forms going blank on StateCompleted: all views now detect
  `StateCompleted` and call `buildForm()` + `form.Init()` to stay live.
- `huh` Select filter mode ESC ejecting the user: ESC is now checked
  before `form.Update` in all form-based views (Select, Analyze, SMD).
- Arrow-key navigation triggering structural rebuilds in mode/band-count
  selects: fixed in both SMDView and InferView by deferring rebuild to
  Enter/Tab confirmation.
- Snapshot memoization: all views now use named snapshot structs to skip
  recomputation when inputs are unchanged.
- `encodeStandardSMD` float boundary overflow: integer bounds check on
  `rounded` prevents wrong-length SMD markings at the 100/1000 boundary.

### Done When:
- ✓ All four views are functional.
- ✓ Navigation between all views works without state corruption.
- ✓ Layout is stable across terminal resize events.
- ✓ `InferView` mode and band count changes do not cause visible form state loss.

---

# Milestone 11 — WASM Module

**Goal:** Expose the core library as a WebAssembly module callable from JavaScript,
so a web application can perform all computation client-side. The server's only
responsibility is serving static files.

### Design Decisions:
- **Build toolchain:** Standard `GOOS=js GOARCH=wasm` is the safe default given
  the library's use of `math`, `strconv`, and `strings`. TinyGo produces smaller
  binaries but has restricted stdlib support; evaluate if binary size becomes a concern.
- **Exported API surface:** At minimum, expose all primary library operations:
  `DecodeBands`, `EncodeBands`, `DecodeSMD`, `EncodeSMD`, `NearestStandard`,
  `SelectStandardResistor`, `InferResistor`, `AnalyzeResistor`.
- **JS interop convention:** All exported functions must accept and return
  JS-compatible types. Errors must surface as JS exceptions or structured return
  objects — panics must not reach the JS caller.
- **Runtime shim:** The standard Go WASM build requires `wasm_exec.js`. This file
  must be versioned alongside the `.wasm` artifact.

### Deliverables:
- `cmd/resistor-wasm/` entry point exporting all core operations as JS-callable functions
- Consistent error handling — no panics crossing the Go/JS boundary
- `make build-wasm` target producing a `.wasm` artifact and accompanying JS shim
- Reference HTML page that exercises each exported function to verify browser compatibility

### Done When:
- All core library operations are callable from browser JavaScript.
- Errors surface cleanly without panics.
- `make build-wasm` is reproducible.
- Reference page loads and operates correctly in a current browser.
