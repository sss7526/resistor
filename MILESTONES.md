# Project Vision

> A Go library that converts between resistor specifications and visual encodings, and performs structured, assumption-aware inference for unknown resistors.

Core principles:
- Deterministic where standards exist
- Explicit assumptions where inference is required
- No UI dependencies in the core library
- Safe defaults for hobbyists

Delivery targets:
- Core library: importable Go package for use in other projects ✓
- CLI: non-interactive command-line toolkit ✓
- TUI: interactive terminal application ✓
- WASM: client-side web toolkit; server responsibility is serving static files only ✓

**All delivery targets complete as of Milestone 11.**

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

# Milestone 11 — WASM Module ✓

**Goal:** Expose the core library as a WebAssembly module callable from JavaScript,
so a web application can perform all computation client-side. The server's only
responsibility is serving static files.

### Design Decisions:
- **Build toolchain:** Standard `GOOS=js GOARCH=wasm` — produces a 3.4 MB binary.
  TinyGo was not needed; stdlib usage (`math`, `strconv`, `strings`, `encoding/json`) is
  fully supported by the standard toolchain.
- **Exported API surface:** All eight primary library operations exposed as JS-callable
  functions on the `resistor` global object.
- **JS interop convention:** Inputs are JSON strings; outputs are plain JS objects
  `{ok: bool, value: any}` or `{ok: false, error: string}`. The `safe()` wrapper
  recovers panics at every function boundary so none reach the JS caller.
- **Runtime shim:** `wasm_exec.js` is copied from `$(go env GOROOT)/lib/wasm/wasm_exec.js`
  at build time and versioned alongside the `.wasm` artifact in `web/`.

### Deliverables:
- ✓ `cmd/resistor-wasm/main.go` — entry point exporting all core operations
- ✓ Consistent error handling — panics recovered, all errors returned as `{ok: false, error}`
- ✓ `make build-wasm` — reproducible; outputs `web/resistor.wasm` + `web/wasm_exec.js`
- ✓ `web/embed.go` + `cmd/resistor-server/` — full web application serving the WASM module

### Done When:
- ✓ All core library operations are callable from browser JavaScript.
- ✓ Errors surface cleanly without panics.
- ✓ `make build-wasm` is reproducible.
- ✓ Reference page loads and operates correctly in a current browser.

---

# Milestone 12 — WASM Binary Size Reduction ✓

**Goal:** Reduce the WASM binary from ~3.4 MB to under 500 KB using TinyGo,
making it practical to serve over the web without a loading penalty.

### Design Decisions:
- **TinyGo 0.41.1** (requires Go ≥1.26) produces a 1.3 MB raw / **430 KB gzip** artifact —
  under the 500 KB target. No changes to `cmd/resistor-wasm/` were required;
  `encoding/json`, `syscall/js`, and all library imports compiled without modification.
- **Build toolchain:** `make build-wasm TINYGO=1` or `make build-wasm-tinygo`.
  TinyGo binary path overridable via `TINYGO_BIN` variable.
- **JS shim:** TinyGo ships its own `wasm_exec.js` at `$(tinygo env TINYGOROOT)/targets/wasm_exec.js`.
  The build target copies it alongside the `.wasm` artifact. It exports the same `Go` class
  as the standard runtime shim, so the server application is fully compatible with both builds.
- **Default build unchanged:** `make build-wasm` still uses the standard Go toolchain
  (3.4 MB / ~950 KB gzip) so the build works without TinyGo installed.

### Constraints:
- The library uses `math`, `strconv`, `strings`, and `encoding/json`. TinyGo supports
  all of these but `encoding/json` reflection support is partial — evaluate whether
  the JSON interop layer in `cmd/resistor-wasm/` needs adjustment.
- If `encoding/json` is incompatible, replace with manual `syscall/js` field
  mapping or a TinyGo-compatible JSON library.
- `make build-wasm` must remain the single build command; add a `TINYGO=1` flag
  or a separate `make build-wasm-tinygo` target.

### Done When:
- ✓ TinyGo build produces a correct `.wasm` artifact that passes the reference page checks.
- ✓ Binary is under 500 KB (gzip) — achieved 430 KB.
- ✓ Fallback standard-Go build still works if TinyGo is not installed.

---

# Milestone 13 — Hot-Path Performance ✓

**Goal:** Replace the two linear-scan hot paths identified by benchmarks with
O(1) or O(log N) alternatives.

### Design Decisions:
- **`NearestStandard` E96/E192:** Two compounding fixes: (1) E48/E96/E192 base tables
  are now pre-computed in `init()` and stored in `eSeriesBase` — eliminates the per-call
  `generateESeries` allocation. (2) Linear scan replaced with `sort.SearchFloat64s`
  (binary search) — O(log N) instead of O(N). `baseValues` simplified to a single map
  lookup. Result: **57 ns/op** (was ~6,500 ns, **~114×** improvement).
- **`EncodeSMD` EIA-96:** `eia96EncodeMap` (96 × 12 = 1,152 entries) built once at
  `init()`. `encodeEIA96` replaced with a single map lookup keyed by
  `roundToSignificant(resistance, 6)` — same rounding used by `decodeEIA96` and
  `NearestStandard`, so float64 key equality holds. Result: **31 ns/op**
  (was ~410 ns, **~13×** improvement).

### Targets:
- **`NearestStandard` E96/E192** (~6.5 µs): currently iterates all values per decade.
  Replace with binary search on the sorted series table.
- **`EncodeSMD` EIA-96** (~410 ns): currently an O(1152) nested loop scan.
  Replace with a pre-built `map[float64]string` lookup keyed by the 96 canonical
  values × 12 multiplier letters.

### Done When:
- ✓ `BenchmarkNearestStandard_E96` improves by ≥ 5× — achieved ~114×.
- ✓ `BenchmarkEncodeSMD_EIA96` improves by ≥ 5× — achieved ~13×.
- ✓ All existing tests continue to pass.

---

# Milestone 14 — E48/E96/E192 SMD Encoding ✓

**Goal:** Close the gap between `NearestStandard` (which supports E48–E192) and
`EncodeSMD` (which only supports 3/4-digit and EIA-96), so high-precision series
values are fully round-trippable through the CLI, TUI, and WASM encode path.

### Design Decisions:
- **`SMDAuto` cascade:** `EncodeSMD` with `SMDAuto` now attempts 3/4-digit first and
  falls through to EIA-96 on failure. `SMDStandard` remains strict (3/4-digit only).
  No new encoding mode was needed; the cascade fits naturally into the existing
  `SMDAuto` semantics.
- **`DecodeSMD` ordering fix:** EIA-96 format check (2 digits + letter) is now tested
  before R-notation. Previously, EIA-96 codes using the 'R' multiplier letter (×0.1,
  e.g. "01R") were misread as R-notation and returned wrong values. The reorder is
  safe because valid R-notation markings ("4R7", "R47") have 'R' in position 0 or 1,
  not position 2, so the EIA-96 pattern check rejects them correctly.
- **CLI and WASM:** both call `EncodeSMD` directly, so the fix is transparent.

### Problem:
`EncodeSMD` with `SMDStandard` or `SMDAuto` rejects values that are valid E96/E192
standard values but not representable in 3/4-digit format. A user who calls
`SelectStandardResistor` with E96 and then tries to encode the result as an SMD
marking gets an error.

### Deliverables:
- ✓ `EncodeSMD` auto-selects EIA-96 for E96/E192 values that are not 3/4-digit
  representable, rather than returning an error.
- ✓ CLI `smd encode` and WASM `encodeSMD` pick up the fix transparently.
- ✓ New tests covering the E96 → SMD round-trip.

### Done When:
- ✓ `SelectStandardResistor` E96 result feeds into `EncodeSMD` without error for
  all 96 base values across all decades.
- ✓ Round-trip `DecodeSMD(EncodeSMD(v, SMDAuto)) == v` holds for all 96 × 12 = 1,152
  E96 values (all base values across all EIA-96 multiplier decades).

---

# Milestone 15 — Web Application ✓

**Goal:** Build a usable, production-ready web application on top of the WASM module,
served by a hardened Go HTTP server, suitable for hobbyists.

### Server Architecture

- **Language/stdlib:** Pure Go, `net/http` only — no third-party web frameworks.
  Use Go 1.22+ `ServeMux` with method+path routing (`GET /`, `GET /health`, etc.).
- **Static file serving:** Embed `web/` via `//go:embed` so the binary is fully
  self-contained. Serve `.wasm` and `.js` with correct `Content-Type` headers.
- **Configuration:** Address and port via flags (`-addr`, default `":8080"`) or
  `RESISTOR_ADDR` env var; flags take precedence. No hardcoded interface or port.
- **Binary location:** `bin/resistor-server`. Makefile target: `build-server`.
  `make build` adds `build-server` to the default build chain.
- **Intended deployment:** Behind a reverse proxy or TLS-termination endpoint
  (nginx, Caddy, Fly.io, etc.). No TLS in the server itself.

### Server Hardening (production defaults, no flags required)

#### Transport & lifecycle
- **Timeouts:** `ReadHeaderTimeout` 2s, `ReadTimeout` 5s, `WriteTimeout` 10s,
  `IdleTimeout` 120s — all set on `http.Server` directly, not overridable at runtime.
- **Request size limit:** `http.MaxBytesReader` wraps every request body (even GET,
  via middleware) to prevent slowloris-style body exhaustion.
- **Panic recovery:** middleware recovers from handler panics and returns `500`
  instead of crashing the process; panic detail is logged server-side, never sent
  to the client.
- **Graceful shutdown:** `os.Signal` handler for `SIGINT`/`SIGTERM`; calls
  `server.Shutdown(ctx)` with a drain timeout so in-flight requests complete.
- **No `Server` header:** strip or replace the default Go server banner to avoid
  version disclosure.

#### Routing & input
- **Method enforcement:** Go 1.22 `ServeMux` pattern `GET /path` — non-matching
  methods receive `405 Method Not Allowed` automatically.
- **No directory listing:** embedded `fs.FS` is served with an explicit file map;
  directory index requests return `404`.
- **Path canonicalisation:** `http.StripPrefix` / `http.RedirectHandler` for any
  trailing-slash normalisation; no open redirects.
- **No user-controlled input reaches the server** — this is a pure static-serving
  app; the health endpoint accepts no parameters. Any query string or body on
  non-health routes is ignored after size-limiting.

#### Security headers (applied by middleware to every response)
| Header | Value |
|--------|-------|
| `Content-Security-Policy` | `default-src 'none'; script-src 'self' 'nonce-{n}'; style-src 'self'; font-src 'self'; img-src 'self' data:; connect-src 'self'; object-src 'none'; base-uri 'self'; form-action 'self'; frame-ancestors 'none'` |
| `X-Content-Type-Options` | `nosniff` |
| `X-Frame-Options` | `DENY` |
| `Referrer-Policy` | `strict-origin-when-cross-origin` |
| `Permissions-Policy` | `camera=(), microphone=(), geolocation=(), payment=(), usb=()` |
| `Cross-Origin-Opener-Policy` | `same-origin` |
| `Cross-Origin-Embedder-Policy` | `require-corp` |
| `Cross-Origin-Resource-Policy` | `same-origin` |

**CSP nonce:** the server generates a 128-bit `crypto/rand` nonce per request,
base64url-encodes it, and injects it into both the `Content-Security-Policy` header
and the HTML template. External scripts (`wasm_exec.js`) are covered by `'self'`.
No `unsafe-inline`, no `unsafe-eval`. `'wasm-unsafe-eval'` is included in `script-src`
— required by Chrome 97+ for `WebAssembly.instantiateStreaming()` regardless of
fetch-based delivery. HTML is rendered via `html/template` (not
`text/template`) which auto-escapes all template values.

#### Health endpoint
`GET /health` returns `200 OK`, `Content-Type: application/json`,
body `{"status":"ok"}`. No auth. Suitable for load-balancer probes.
All security headers still applied.

### Client-Side Security (XSS and DOM attack surface)

- **`html/template` rendering:** all server-injected values (nonce, version, etc.)
  pass through Go's `html/template` context-aware escaper — safe in HTML, attribute,
  JS, and URL contexts automatically.
- **No `innerHTML`, no `document.write`, no `eval`:** all DOM updates use
  `textContent`, `createElement`, and `appendChild`. No string-to-DOM paths exist.
- **No `dangerouslySetInnerHTML` or equivalent:** WASM results are plain JS objects
  (numbers, strings, arrays); they are displayed via `textContent` only, never
  interpolated into markup strings.
- **Trusted Types** (where supported): `require-trusted-types-for 'script'` added
  to the CSP so browsers that support Trusted Types enforce the no-innerHTML rule
  at the platform level.
- **No third-party scripts, fonts, or stylesheets:** all assets are same-origin
  and embedded in the binary. `default-src 'none'` in the CSP means any accidental
  external load is blocked.
- **No `postMessage` cross-origin:** the app does not use iframes or `postMessage`.
  `frame-ancestors 'none'` prevents the page being framed by any origin.
- **Form inputs:** all `<input>` elements have explicit `type` attributes to prevent
  type-confusion. Numeric inputs use `type="number"` with `min`/`max`/`step`
  constraints validated client-side by the browser before being passed to WASM.
  WASM itself validates all inputs independently (defense in depth).
- **No cookies, no localStorage, no sessionStorage:** the app is stateless;
  no user data is persisted anywhere.
- **`COEP: require-corp` + `COOP: same-origin`:** enables `SharedArrayBuffer`
  isolation if ever needed, and prevents Spectre-class cross-origin data leaks.

### Views:
| View | Description |
|---|---|
| ✓ Select Resistor | Enter resistance + series, get snapped standard value and color band diagram |
| ✓ Decode Bands | Select band colors; read resistance and tolerance |
| ✓ Encode Bands | Enter resistance + tolerance; get color band diagram |
| ✓ SMD Tools | Decode or encode SMD markings (auto / standard / EIA-96 modes) |
| ✓ Infer Resistor | Enter physical observations; read inferred properties and confidence score |
| ✓ Analyze Resistor | Enter electrical conditions; read power dissipation, derating, and warnings |

### Done When:
- ✓ `make build-server` produces `bin/resistor-server`.
- ✓ Server starts with `./bin/resistor-server -addr :8080` and serves the app.
- ✓ `-addr` flag (or `RESISTOR_ADDR` env) controls listen address; no defaults
  are hardcoded in source.
- ✓ All views are functional end-to-end via WASM.
- ✓ Color band diagram renders correct colors for any valid input.
- ✓ All security headers from the table above are present on every response.
- ✓ CSP nonce rotates per request; no two responses share a nonce.
- ✓ `curl -v /health` returns `200` with `{"status":"ok"}`.
- ✓ Graceful shutdown on `SIGINT`/`SIGTERM` — in-flight requests complete.
- ✓ `Server` response header absent or replaced (no Go version disclosure).
- ✓ No `innerHTML`, `eval`, `document.write`, or `new Function` anywhere in JS.
- ✓ No third-party resources loaded; `default-src 'none'` CSP blocks any accidental
  external load.
- ✓ Works in current Chrome, Firefox, and Safari without a build step on the client.

---

# Milestone 16 — Version Release

**Goal:** Cut a tagged `v0.1.0` release so the `go get` path in the README resolves
and the module is addressable by version from external projects.

### Deliverables:
- `v0.1.0` git tag pushed to the remote
- GitHub release created with release notes summarising the API surface
- `go get github.com/sss7526/resistor@v0.1.0` resolves correctly via the module proxy

### Done When:
- Tag exists on the remote and is visible via `go list -m github.com/sss7526/resistor@v0.1.0`.
- Release notes cover all public API entry points.
