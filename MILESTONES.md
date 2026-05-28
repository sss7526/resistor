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
- ✓ Footer CTA links to project source (GitHub) and documentation (pkg.go.dev) with
  theme-consistent SVG icons; no external resources loaded.

---

# Milestone 16 — Containerised Deployment ✓

**Goal:** Package the web server into a minimal, production-ready container image
and provide a complete Caddy-fronted deployment configuration.

### Deliverables

- **`Dockerfile`** — multi-stage build; standard Go WASM by default, TinyGo opt-in
  via `--build-arg WASM=tinygo`. Final stage is `scratch`; binary is statically
  linked (`CGO_ENABLED=0 -ldflags="-s -w -trimpath"`).
- **`docker-compose.yml`** — `resistor` + `caddy:2-alpine` services on an internal
  network; only Caddy exposes ports 80/443; Caddy volumes for cert persistence.
- **`Caddyfile`** — hostname from `$RESISTOR_HOST` env var; automatic HTTPS via
  Let's Encrypt; reverse-proxies to `resistor:8080`.
- **README** — deployment section covering Docker build, Compose, and Caddyfile
  customisation.

### Container Constraints

- Final image: `scratch` (no shell, no OS, binary only).
- Run as UID 10001 (non-root, no `/etc/passwd` required in scratch).
- No healthcheck command in image (scratch has no tooling); health checks delegated
  to Caddy active probing or orchestrator HTTP checks against `/health`.
- Single binary embeds all static assets — no volume mounts needed at runtime.

### Done When:
- ✓ `docker build .` produces a working image; `docker run -p 8080:8080 <image>` serves the UI.
- ✓ `docker build --build-arg WASM=tinygo .` produces a smaller image with TinyGo WASM.
- ✓ `docker compose up` starts `resistor` + `caddy`; UI reachable via Caddy on port 80.
- ✓ `Caddyfile` `$RESISTOR_HOST` env var controls the public hostname.
- ✓ README deployment section is complete and accurate.

---

# Milestone 17 — Version Release & Automated Maintenance

**Goal:** Ship `v0.1.0`, then wire up automation so the project self-maintains
(dependency updates, future releases, binary publishing, Docker Hub) without
manual intervention between active development milestones.

---

### One-Time Human Setup (do these once, never again)

These steps cannot be automated — they require human hands on external services.
Complete them before implementing the workflows below.

#### 1. Docker Hub access token

1. Log in at [hub.docker.com](https://hub.docker.com).
2. Account Settings → Security → Access Tokens → **New Access Token**.
3. Name it `resistor-goreleaser`, scope: **Read & Write**.
4. Copy the token — it is shown only once.

#### 2. GitHub Actions secrets

Go to the repo on GitHub → Settings → Secrets and variables → Actions →
**New repository secret**. Add both:

| Secret name | Value |
|---|---|
| `DOCKERHUB_USERNAME` | your Docker Hub username |
| `DOCKERHUB_TOKEN` | the access token from step 1 |

#### 3. GitHub Actions permissions

Repo → Settings → Actions → General → Workflow permissions →
set to **Read and write permissions** and enable
**Allow GitHub Actions to create and approve pull requests**.
This is required for release-please to open its Release PR.

That is the complete manual setup. Everything else is automated from this point.

---

### Part 1 — Initial Release

Push the first tagged release so the module is addressable via `go get` and
indexed by `pkg.go.dev`.

The tag and GitHub Release are created automatically by release-please (Part 2)
when its first Release PR is merged. Do not tag manually.

**Deliverables:**
- `v0.1.0` git tag created by release-please on first Release PR merge.
- GitHub Release published automatically with generated changelog.
- `go get github.com/sss7526/resistor@v0.1.0` resolves via the module proxy.

---

### Part 2 — CI Workflow

A dedicated CI workflow runs on every push and every pull request to `main`,
keeping the badge green and gating all merges.

**File:** `.github/workflows/ci.yml`

```yaml
on:
  push:
  pull_request:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod
      - run: make test-all
      - run: make smoke
```

This is also what the CI badge in the README points to.

---

### Part 3 — Automated Releases (release-please)

**How it works:**

`release-please` reads conventional commit messages (`feat:`, `fix:`, `chore:`,
`docs:`, etc.) — already the style used in this repo — and maintains a rolling
"Release PR" on `main`.

- `fix:` commits → **patch** bump (`v0.1.0` → `v0.1.1`)
- `feat:` commits → **minor** bump (`v0.1.0` → `v0.2.0`)
- `feat!:` or `BREAKING CHANGE:` footer → **major** bump (`v0.1.0` → `v1.0.0`)

**To cut a release:** merge the Release PR that release-please keeps open.
That merge triggers the action to create the git tag and publish the GitHub
Release with a generated `CHANGELOG.md` entry. No manual tagging. No manual
release notes. Between milestones the Release PR just accumulates commits and
sits open until you decide to ship.

**Files:**
- `.github/workflows/release-please.yml` — triggers on push to `main`; uses
  `googleapis/release-please-action`; configured for `go` release type;
  creates a **draft** GitHub Release (GoReleaser publishes it with artifacts).
- `release-please-config.json` — repo root; tells release-please the release
  type and package path.
- `.release-please-manifest.json` — repo root; tracks the current version;
  release-please updates this file automatically on each release.

---

### Part 4 — Automated Dependency Maintenance (Dependabot + auto-merge)

Dependabot opens PRs automatically when new versions of dependencies are
available. A companion workflow auto-merges those PRs when CI passes.

**Dependabot config** (`.github/dependabot.yml`) monitors three ecosystems weekly:

| Ecosystem | What it updates |
|---|---|
| `gomod` | `go.mod` / `go.sum` — library and tool dependencies |
| `github-actions` | Action versions pinned in `.github/workflows/` |
| `docker` | Base image tags in `Dockerfile` (`golang:1.26`, `tinygo/tinygo:0.41.0`) |

**Auto-merge workflow** (`.github/workflows/dependabot-automerge.yml`):
- Triggers on `pull_request` events from the `dependabot[bot]` actor.
- Runs the full CI suite (`make test-all`).
- **Minor or patch** updates: enables GitHub auto-merge if CI passes.
- **Major** version bumps: left open for manual review.
- Uses only `GITHUB_TOKEN` — no extra secrets needed.

Security patches and routine updates merge themselves. The only Dependabot PRs
that need attention are major version bumps, which happen rarely.

---

### Part 5 — Version Wiring (CLI, TUI, Server)

**Current state:**

| Binary | Version var | Injected? | Exposed via |
|---|---|---|---|
| `resistor-cli` | `cmd/resistor-cli/cmd.version` | ✓ ldflags | `resistor-cli version` |
| `resistor-tui` | `cmd/resistor-tui/app.version` | ✓ ldflags | title bar |
| `resistor-server` | none | ✗ | nothing |

**Makefile — dynamic VERSION from git:**

```makefile
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
```

Resolves to `v0.1.1`, `v0.2.0-3-gabcdef`, etc. automatically as release-please
creates tags. No Makefile edits ever needed between releases.

**Server version support:**

Add `var version = "dev"` to `cmd/resistor-server/main.go` and expose it:
1. `-version` flag — prints version and exits (mirrors CLI behaviour).
2. `/health` response — extend to `{"status":"ok","version":"v0.1.0"}` so the
   running version is queryable without exec.

Add to `build-server` Makefile target:
```makefile
-ldflags "-X 'main.version=$(VERSION)' -s -w"
```

**Frontend version display:**

On `DOMContentLoaded`, `app.js` fetches `/health` and renders the version string
in the footer alongside the existing GitHub and pkg.go.dev links:

```js
fetch('/health')
  .then(r => r.json())
  .then(d => {
    const el = document.getElementById('app-version')
    if (el && d.version) el.textContent = d.version
  })
  .catch(() => {})   // silently ignore — version display is non-critical
```

`app.html` footer gets a `<span id="app-version"></span>` placeholder. Styled
muted/small so it doesn't compete with the CTA links. Shows `v0.1.0`, `dev`,
or nothing on fetch failure.

**Local Docker builds (for dev/testing):**

`.git` is not in the Docker build context so `git describe` cannot run inside
the container. Makefile targets pass version explicitly:

```makefile
docker-build:
	docker build --build-arg WASM=go --build-arg VERSION=$(VERSION) \
	  -t resistor-server:$(VERSION) .

docker-build-tinygo:
	docker build --build-arg WASM=tinygo --build-arg VERSION=$(VERSION) \
	  -t resistor-server:$(VERSION)-tinygo .
```

`Dockerfile` declares `ARG VERSION=dev` and threads it through the `builder`
stage via `-ldflags "-X 'main.version=${VERSION}'"`.

---

### Part 6 — GoReleaser (Binary Publishing + Docker Hub)

GoReleaser fires on every git tag (created by release-please) and:
1. Builds cross-platform binaries for `resistor-cli` and `resistor-tui`.
2. Uploads them as assets to the GitHub Release that release-please created.
3. Builds both Docker image variants and pushes them to Docker Hub.

**Pipeline flow (nothing manual after initial setup):**

```
merge Release PR
  → release-please creates tag + draft GitHub Release
    → GoReleaser workflow fires on tag
        → builds CLI + TUI binaries (linux/darwin/windows × amd64/arm64)
        → uploads binaries + checksums to GitHub Release, publishes it
        → builds resistor-server:vX.Y.Z (Go WASM)
        → builds resistor-server:vX.Y.Z-tinygo + :latest (TinyGo WASM)
        → pushes both images to Docker Hub
```

**Docker Hub repository:** auto-created on first push — no manual repo creation
needed. The repo will appear at `hub.docker.com/r/sss7526/resistor` after the
first GoReleaser run. The Hub description/README can be updated manually any
time after, but is not required for the images to work.

**`.goreleaser.yml` key sections:**

```yaml
before:
  hooks:
    - make build-wasm   # standard Go WASM must exist before server build

builds:
  - id: resistor-cli
    main: ./cmd/resistor-cli
    ldflags: ["-X 'github.com/sss7526/resistor/cmd/resistor-cli/cmd.version={{.Version}}'"]
    goos: [linux, darwin, windows]
    goarch: [amd64, arm64]

  - id: resistor-tui
    main: ./cmd/resistor-tui
    ldflags: ["-X 'github.com/sss7526/resistor/cmd/resistor-tui/app.version={{.Version}}'"]
    goos: [linux, darwin, windows]
    goarch: [amd64, arm64]

  - id: resistor-server
    main: ./cmd/resistor-server
    ldflags: ["-X 'main.version={{.Version}}' -s -w"]
    goos: [linux]
    goarch: [amd64, arm64]

dockers:
  - image_templates: ["sss7526/resistor:{{.Version}}"]
    build_flag_templates:
      - "--build-arg=WASM=go"
      - "--build-arg=VERSION={{.Version}}"

  - image_templates: ["sss7526/resistor:{{.Version}}-tinygo", "sss7526/resistor:latest"]
    build_flag_templates:
      - "--build-arg=WASM=tinygo"
      - "--build-arg=VERSION={{.Version}}"

release:
  draft: false        # release-please already created the draft; GoReleaser publishes it
  use_existing: true  # attach artifacts to release-please's release, don't create a new one
```

**GoReleaser workflow** (`.github/workflows/goreleaser.yml`):
- Triggers on `push` to tags matching `v*`.
- Logs in to Docker Hub using `DOCKERHUB_USERNAME` and `DOCKERHUB_TOKEN` secrets
  (set up in one-time setup above).
- Runs `goreleaser release --clean`.

**Files:**
- `.goreleaser.yml` — repo root; GoReleaser config.
- `.github/workflows/goreleaser.yml` — triggers on tag push.

---

### Part 7 — EC2 Deployment (one-time bootstrap, then hands-off)

After GoReleaser pushes `:latest` to Docker Hub, Watchtower on the EC2 instance
picks it up automatically within 5 minutes and restarts the container. The only
human steps are the initial instance setup — never repeated.

#### One-time AWS setup

1. **Launch EC2 instance** — Ubuntu 24.04 LTS, `t3.micro` or larger.
2. **Allocate an Elastic IP** and associate it with the instance. Without this
   the public IP changes on every reboot, breaking DNS and TLS cert renewal.
3. **Security group** — inbound rules:

   | Port | Protocol | Source |
   |---|---|---|
   | 22 | TCP | your IP only |
   | 80 | TCP | 0.0.0.0/0 |
   | 443 | TCP | 0.0.0.0/0 |
   | 443 | UDP | 0.0.0.0/0 (HTTP/3) |

4. **Route 53** — create an `A` record pointing your domain → the Elastic IP.
   DNS must resolve to the instance before first startup so Caddy can obtain
   the Let's Encrypt TLS certificate.

#### One-time instance bootstrap (single SSH session)

```bash
# Install Docker (includes Compose plugin)
curl -fsSL https://get.docker.com | sh
sudo usermod -aG docker ubuntu
newgrp docker   # apply group without logout

# Create app directory
mkdir ~/resistor && cd ~/resistor

# Copy docker-compose.yml and Caddyfile from the repo, then start
RESISTOR_HOST=your-domain.com docker compose up -d
```

Caddy provisions the TLS certificate automatically on first startup. Watchtower
begins polling Docker Hub immediately. Close the SSH session — nothing more to do.

#### How future deploys work

```
merge Release PR
  → GoReleaser pushes sss7526/resistor:latest (TinyGo image) to Docker Hub
    → Watchtower detects new digest within 5 minutes
      → pulls new image, recreates resistor container
        → site is running the new version, zero manual action
```

Watchtower logs on the instance (`docker logs watchtower`) show every pull and
restart with timestamps, providing a lightweight audit trail.

#### What never changes after bootstrap
- No SSH for deploys
- No manual `docker pull` or `docker compose up`
- No cert renewal — Caddy handles it automatically
- No version tracking — Watchtower follows `:latest`

---

### Done When:
- All one-time setup steps completed (Docker Hub token, two GitHub secrets,
  Actions permissions).
- `ci.yml` workflow passes on a push to `main`; README CI badge is green.
- `release-please.yml` is present; on first merge to `main` it opens a Release PR.
- Merging the Release PR creates the `v0.1.0` tag and draft GitHub Release.
- GoReleaser fires on the tag, uploads CLI + TUI binaries for all platforms,
  builds and pushes both Docker images to Docker Hub, publishes the release.
- `go get github.com/sss7526/resistor@v0.1.0` resolves via the module proxy.
- `sss7526/resistor:v0.1.0` and `sss7526/resistor:v0.1.0-tinygo` are live on
  Docker Hub — no manual `docker push` performed.
- `.github/dependabot.yml` is present; Dependabot opens its first batch of PRs
  within 24 hours.
- `dependabot-automerge.yml` auto-merges a patch-level Dependabot PR after CI
  passes (verified on first real Dependabot PR).
- All three binaries report the correct release tag: `resistor-cli version`,
  `resistor-tui` title bar, `resistor-server -version` and `/health`.
- Web UI footer displays the version fetched from `/health` on page load.
- No manual tagging, version editing, binary uploading, or docker pushing is
  required for any future release — merging the Release PR is the only action.
- EC2 instance bootstrapped: Elastic IP assigned, Route 53 A record resolves,
  Caddy TLS cert issued, `docker compose up -d` running with Watchtower.
- A test release confirms end-to-end: merge Release PR → new image on Docker Hub
  → Watchtower restarts container on EC2 within 5 minutes, no SSH required.

---

# Milestone 18 — Circuit-Level Analysis

**Goal:** Extend the library and web UI beyond single-component lookup into practical
circuit calculations that hobbyists reach for every session at the bench.

### Motivation

The single-resistor model covers identification and basic safety checking well.
The next most common real-world tasks are all multi-component or multi-condition:
figuring out a voltage divider ratio, combining resistors from a limited parts bin,
sizing an LED current-limiter, and reasoning about RC timing. These share a clean
boundary — they are all resistor-centric circuits with deterministic closed-form
solutions — making them a natural and cohesive second layer.

### Library additions (`github.com/sss7526/resistor`)

#### Series / Parallel

```go
// SeriesResistance returns the total resistance of resistors in series.
func SeriesResistance(values ...float64) (float64, error)

// ParallelResistance returns the equivalent resistance of resistors in parallel.
func ParallelResistance(values ...float64) (float64, error)
```

Both validate that all inputs are positive and that at least one value is provided.
`ParallelResistance` guards against division by zero on a zero-valued input.

#### Voltage Divider

```go
type VoltageDividerResult struct {
    OutputVoltage float64  // Vout = Vin × R2 / (R1 + R2)
    CurrentMA     float64  // quiescent current in mA
    R1            float64  // top resistor
    R2            float64  // bottom resistor
}

// VoltageDivider computes output voltage and quiescent current for a resistive
// voltage divider given supply voltage, R1 (top), and R2 (bottom).
func VoltageDivider(vin, r1, r2 float64) (VoltageDividerResult, error)

type DividerRequest struct {
    Vin          float64   // supply voltage
    TargetVout   float64   // desired output voltage
    Series       ESeries   // preferred E-series for resistor selection
    MaxCurrentMA float64   // quiescent current budget (0 = unconstrained)
}

// FindDividerPair selects the closest standard E-series R1/R2 pair that achieves
// the target output voltage within the current budget.
func FindDividerPair(req DividerRequest) (VoltageDividerResult, error)
```

#### LED Current Limiter

```go
type LEDCurrentLimiterResult struct {
    Resistor         float64  // calculated series resistance in Ω
    StandardResistor float64  // nearest E-series value
    ActualCurrentMA  float64  // current with standard resistor
    PowerDissipation float64  // power in the resistor (W)
}

// LEDCurrentLimiter calculates the series resistor needed to drive an LED at a
// target current given supply voltage and LED forward voltage.
func LEDCurrentLimiter(vsupply, vforward, targetCurrentMA float64, series ESeries) (LEDCurrentLimiterResult, error)
```

#### RC Time Constant

```go
type RCResult struct {
    TauSeconds    float64  // τ = R × C
    FrequencyHz   float64  // -3 dB frequency = 1 / (2π × τ)
    Rise10_90MS   float64  // 10–90 % rise time ≈ 2.197 × τ (ms)
}

// RCTimeConstant computes timing and frequency characteristics for a series RC
// circuit. R in ohms, C in farads.
func RCTimeConstant(r, c float64) (RCResult, error)
```

### Web UI additions

Two new tabs in the existing 6-tab layout:

| Tab | Description |
|---|---|
| Networks | Series/parallel combination; enter 2–8 values, get equivalent resistance and E-series snap |
| Circuits | Voltage divider (compute or find pair), LED current limiter, RC time constant |

Existing tabs and API surface unchanged.

### WASM additions

Three new exported functions on the `resistor` global object following the existing
`{ok, value}` / `{ok, error}` envelope:

```js
resistor.voltageDivider('{"vin":5,"r1":10000,"r2":4700}')
resistor.findDividerPair('{"vin":5,"targetVout":3.3,"series":"E24","maxCurrentMA":1}')
resistor.ledCurrentLimiter('{"vsupply":5,"vforward":2.1,"targetCurrentMA":20,"series":"E24"}')
resistor.rcTimeConstant('{"r":10000,"c":0.0000001}')
resistor.seriesResistance('[1000,2200,4700]')
resistor.parallelResistance('[1000,1000]')
```

### Done When:
- All new library functions have unit tests covering normal cases, edge cases
  (single input, zero, negative), and round-trip consistency with existing API.
- `FindDividerPair` result always uses valid E-series values verifiable by
  `NearestStandard`.
- WASM exports functional; existing WASM tests unaffected.
- All two new UI tabs operational end-to-end.
- `make test-all` and `make smoke` pass.
