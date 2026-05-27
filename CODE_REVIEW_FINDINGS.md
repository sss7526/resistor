# Code Review Findings

**This file is intentionally temporary.** It tracks open defects found during a code review of
the initial implementation. Delete it once every item below is resolved, verified, and the CI
gates pass clean. Do not let it drift — close items as they are fixed, not in a batch at the end.

---

## CI Gates

Before committing any fix, run both gates and confirm they pass:

```
make test-all
make smoke
```

`make test-all` runs unit tests and CLI integration tests. `make smoke` exercises the built binary
end-to-end. A fix that breaks either gate must be corrected before the commit lands.

---

## Findings

Ranked by severity. Each entry states the defect, how to fix it, and whether tests need to be
added or revised.

---

### 1. `--json` exits 0 on command failure

**File:** `internal/cli/json.go:43`

**Defect:** `Respond()` calls `OutputJSONError` and then returns `nil` when `jsonOutput=true` and
`err!=nil`. Cobra sees a `nil` return from `RunE` and exits with code 0. Every command failure
under `--json` is invisible to shell scripts and CI pipelines that check `$?`.

**Fix:** Return the error after printing it, not nil:

```go
if jsonOutput {
    _ = OutputJSONError(err)
    return err
}
```

**Tests:**
- No existing test covers exit codes. Add a case to `cmd/resistor-cli/integration_test.go` that
  invokes a failing command with `--json` and asserts the exit code is non-zero. A straightforward
  way is to run the command via `exec.Command`, check `cmd.ProcessState.ExitCode()`, and also
  verify the stdout contains `"success":false`.

---

### 2. `NearestStandard(RoundUp)` returns wrong value at decade boundary

**File:** `eseries.go:255`

**Defect:** When the normalized input exceeds the largest entry in the base decade table (e.g.,
9.15 against E24 whose max is 9.1), the `RoundUp` loop finds no candidate and falls through to the
shared return at line 269. At that point `best` still holds its initialization value of `base[0]`
(1.0), so the function returns `1.0 × 10^exponent` — the bottom of the current decade — instead
of `1.0 × 10^(exponent+1)`.

Example: `NearestStandard(9.15, E24, RoundUp)` returns `1.0` instead of `10.0`.

**Fix:** After the loop, handle the fallthrough for `RoundUp` explicitly:

```go
// After the loop:
if mode == RoundUp {
    // No candidate in this decade was >= normalized; step to first value of next decade.
    result := base[0] * math.Pow(10, exponent+1)
    return roundToSignificant(result, 6), nil
}
```

**Tests:**
- `TestNearestStandard_RoundUp` in `eseries_test.go` has three cases, none of which hit a decade
  boundary. Add at least one case that does:
  ```
  {name: "E24 9.15Ω → 10Ω (decade wrap)", input: 9.15, series: E24, expected: 10.0}
  ```
  Also add a kΩ-scale variant to confirm the exponent is handled correctly:
  ```
  {name: "E24 9150Ω → 10000Ω (decade wrap)", input: 9150, series: E24, expected: 10000.0}
  ```

---

### 3. Ohm's Law consistency check is dead

**File:** `analysis.go:94`

**Defect:** When both `AppliedVoltage` and `AppliedCurrent` are provided, the check fires only
when `|I×R − V| > 1e6` (one million volts). For any physically realizable input the discrepancy
will never reach that threshold. The `WarningCaution` for inconsistent V/I is never emitted.

Example: V=10, I=0.001, R=100Ω → discrepancy is 9.9; threshold is 1,000,000.

**Fix:** Replace the absolute threshold with a relative one. A tolerance of 1% of the expected
voltage is reasonable for this warning:

```go
if math.Abs(expectedV-V)/math.Max(math.Abs(V), 1e-12) > 0.01 {
```

Alternatively use an absolute threshold commensurate with realistic voltages (e.g., `0.01` volts
for a 1% tolerance at low voltages, combined with the relative check).

**Tests:**
- No existing test exercises the both-V-and-I path and checks for the inconsistency warning.
  Add a case to `TestAnalyzeResistor` that provides V and I that disagree by more than 1%, and
  asserts a `WarningCaution` appears in `report.Warnings`.
- Add a complementary case where V and I are consistent (within 1%) and asserts no inconsistency
  warning is emitted.

---

### 4. `InferResistor` silently discards decode errors, always returns nil error

**File:** `inference.go:113` and `inference.go:126`

**Defect:** When `DecodeBands` or `DecodeSMD` returns an error, the `if err == nil` guard skips
the block silently. The error is not returned, not logged, and not recorded in `Meta.Assumptions`.
`InferResistor` unconditionally returns `nil` as its error. Callers cannot distinguish "invalid
input that failed to decode" from "not enough information to infer."

**Fix:** When a deterministic decode was attempted and failed, surface the failure. There are two
reasonable approaches:

**Option A — record in assumptions** (preserves the nil-error contract, adds visibility):
```go
spec, err := DecodeBands(obs.Bands)
if err == nil {
    result.Spec = spec
    assumptions = append(assumptions, "Resistance and tolerance determined from color bands")
    confidence = 1 - (1-confidence)*(1-1.0)
} else {
    assumptions = append(assumptions, "Color band decode failed: "+err.Error())
}
```

**Option B — return the error** (stricter contract, breaking change):
```go
spec, err := DecodeBands(obs.Bands)
if err != nil {
    return InferenceResult{}, fmt.Errorf("band decode: %w", err)
}
```

Option A is the lower-disruption choice and fits the existing design (inference is best-effort).
Option B is cleaner for callers that want hard failures. Choose based on the intended contract
for the public API before making this change.

**Tests:**
- No existing test passes invalid band colors to `InferResistor`. Add a case with well-formed
  band count (4) but invalid colors and assert the failure is visible — either via an error return
  or via a recorded assumption — depending on which option is chosen.
- Add a parallel case for an invalid SMD marking with a valid format but undecodable content.

---

### 5. Exported `TempCoeffPPM` map is unused and diverges from `TempCoeffValue`

**File:** `mappings.go:90`

**Defect:** `var TempCoeffPPM` is an exported map with a doc comment describing it as the
authoritative color-to-PPM mapping. However, all internal code — `DecodeBands` (line 167) and
CLI helpers — reads `TempCoeffValue` exclusively. `TempCoeffPPM` is never read anywhere in the
codebase. Any caller or future maintainer who edits `TempCoeffPPM` (e.g., adding a new entry)
will see no effect on behavior. The two maps will silently diverge.

The name collision with the struct field `ResistorSpec.TempCoeffPPM` adds further confusion.

**Fix:** Remove `var TempCoeffPPM`. It is dead exported surface. If the intent was to provide a
public API for the map, rename `TempCoeffValue` to `TempCoeffPPM` and update all internal
references, or explicitly re-export it:

```go
// TempCoeffPPM is the public name for the 6th-band color-to-PPM mapping.
var TempCoeffPPM = TempCoeffValue
```

The second approach creates a true alias rather than a duplicate.

**Tests:**
- No new tests are required. This is a dead symbol removal. Confirm that deleting or aliasing
  `TempCoeffPPM` does not break compilation — `go build ./...` and `make test-all` are sufficient.

---

### 6. `findEIA96Multiplier` is dead code with a broken implementation

**File:** `smd.go:255`

**Defect:** `findEIA96Multiplier` is defined but never called anywhere in the codebase. Its
tolerance condition `math.Abs(v-mult) < 1e9` is also broken: since the largest multiplier in
`eia96Multipliers` is `1e8`, the condition is always true regardless of `mult`, causing the
function to return whichever map entry Go's non-deterministic iteration visits first.

**Fix:** Delete the function. There are no call sites to update.

**Tests:**
- No tests reference this function. No test changes required.

---

### 7. `ErrInvalidBandCount` message omits 6-band support

**File:** `bands.go:22`

**Defect:** The error message reads `"invalid number of bands (must be 4 or 5)"` but `DecodeBands`
accepts 4, 5, and 6 bands. The doc comment on `ErrInvalidBandCount` and the `DecodeBands` function
comment also only mention 4 and 5. A caller who receives this error for a 3-band or 7-band input
is misinformed about what counts are valid.

**Fix:** Update the error message and both doc comments:

```go
var ErrInvalidBandCount = errors.New("invalid number of bands (must be 4, 5, or 6)")
```

Update the doc comment on `ErrInvalidBandCount` and the `DecodeBands` function header to mention
6-band support.

**Tests:**
- `TestDecodeBands_InvalidCases` in `bands_test.go` tests that invalid counts return an error, but
  does not assert the message text. No test changes required for the fix to be safe, but consider
  adding a `require.ErrorIs` assertion to pin the sentinel value so future message changes don't
  go unnoticed.

---

### 8. `strings.Title` deprecated since Go 1.18

**File:** `internal/cli/format.go:16`

**Defect:** `strings.Title(string(b))` is deprecated as of Go 1.18. The module declares
`go 1.26.2`. The replacement package `golang.org/x/text/cases` is already present in `go.mod` as
an indirect dependency, so no new dependency is required.

**Fix:**

```go
import "golang.org/x/text/cases"
import "golang.org/x/text/language"

var titler = cases.Title(language.Und)

func PrintBands(bands []resistor.Color) {
    fmt.Println("Bands:")
    for _, b := range bands {
        fmt.Printf("  %s\n", titler.String(string(b)))
    }
}
```

The `language.Und` (undetermined) tagger is appropriate here since band color names are ASCII
identifiers, not natural-language text.

**Tests:**
- No behavioral change for ASCII color names. No test changes required.

---

## Lint Findings (golangci-lint v2.12.2)

Run with: `make lint`

Findings #6 and #8 above are confirmed independently by the `unused` and `staticcheck` linters.
The items below are new findings surfaced by lint that are not covered by the code review section.

---

### 9. `smd_test.go:135` — integer overflow conversion and nonsensical test names

**File:** `smd_test.go:135`
**Linter:** gosec G115

**Defect:** `string(rune(int(val)))` converts a `float64` test value to `int`, then to `rune`, then
to a string to form a test name. gosec flags the `int → rune` conversion as a potential overflow
(on 32-bit platforms `int` is 32 bits, same as `rune`, but the conversion is still flagged). Beyond
the linter warning, the resulting test names are meaningless Unicode codepoints: `val=100` produces
`"d"`, `val=1000` produces `"Ϩ"`, `val=4990` produces a random Unicode character. These names are
useless for diagnosing failures.

**Fix:** Replace with a proper numeric string for the test name:

```go
t.Run(fmt.Sprintf("EIA96 encode %.0f", val), func(t *testing.T) {
```

**Tests:** This is the test itself. No additional test needed.

---

### 10. `integration_test.go:27,42` — gosec G204 false positives in test code

**File:** `cmd/resistor-cli/integration_test.go:27` and `:42`
**Linter:** gosec G204

**Defect:** gosec flags `exec.Command` calls where the binary path or arguments come from
variables. Both call sites are intentional: one builds the CLI binary, the other invokes it under
test. Neither is a security issue — integration tests necessarily spawn subprocesses with
variable paths.

**Fix:** Suppress with `//nolint:gosec` at each call site, with a brief reason:

```go
build := exec.Command("go", "build", "-ldflags", ldflags, "-o", binaryPath, ".") //nolint:gosec // intentional subprocess in test
...
cmd := exec.Command(binary, args...) //nolint:gosec // intentional subprocess in test
```

**Tests:** No test changes needed beyond the suppression comment.

---

### 11. `helpers.go:31` — De Morgan's law (QF1001)

**File:** `cmd/resistor-cli/cmd/helpers.go:31`
**Linter:** staticcheck QF1001

**Defect:** The condition `!(digitOK || multOK || tolOK || tempOK)` is logically equivalent to
`!digitOK && !multOK && !tolOK && !tempOK` but is harder to read. staticcheck suggests applying
De Morgan's law.

**Fix:**

```go
if !digitOK && !multOK && !tolOK && !tempOK {
```

**Tests:** Logic is unchanged. No test changes required.

---

### 12. `WriteString(fmt.Sprintf(...))` pattern (QF1012)

**Files:** `cmd/resistor-tui/app/common.go:77,89` and `cmd/resistor-tui/app/select.go:207`
**Linter:** staticcheck QF1012

**Defect:** Three call sites use `b.WriteString(fmt.Sprintf(...))` where `fmt.Fprintf(&b, ...)`
is both shorter and avoids an intermediate string allocation.

**Fix:** Replace each occurrence:

```go
// Before
b.WriteString(fmt.Sprintf("  %s\n", c))
// After
fmt.Fprintf(&b, "  %s\n", c)
```

Apply the same pattern at all three locations.

**Tests:** Output is identical. No test changes required.

---

### 13. TUI stub symbols flagged as unused

**Files:** `cmd/resistor-tui/app/model.go:7–15`, `cmd/resistor-tui/app/common.go:74,94`,
`cmd/resistor-tui/app/version.go:3`
**Linter:** unused

**Defect:** Nine symbols are unused because the TUI's Analyze and SMD views are not yet
implemented:

- `type viewState` and constants `viewMenu`, `viewSelect`, `viewInfer`, `viewAnalyze`, `viewSMD`,
  `viewQuit` in `model.go` — view routing enum defined but not wired to any real views
- `func renderBands` and `func renderWarnings` in `common.go` — helpers intended for Analyze/SMD
  views that don't exist yet
- `var version` in `version.go` — injected via `-ldflags` at build time but not referenced in any
  Go code path

**Fix:** These will clear naturally when the TUI Analyze and SMD views are implemented (tracked in
MILESTONES.md M10). Until then the options are:

- Leave `make lint` failing on these — acceptable while the TUI is a known work in progress
- Suppress individually with `//nolint:unused` — adds noise to stub code
- Wire `version` into the TUI's help or about text now, which is a trivial one-liner and removes
  one warning immediately

The `var version` unused warning is the easiest to fix independently of the view work:

```go
// version.go — ensure the variable is referenced somewhere, e.g. in model.go:
func (m AppModel) version() string { return version }
```

Or simply use it in an existing string displayed by the TUI.

**Tests:** No test changes required.

---

## Resolution Checklist

| # | File | Status |
|---|------|--------|
| 1 | `internal/cli/json.go:43` — `--json` exits 0 on error | open |
| 2 | `eseries.go:255` — `RoundUp` decade boundary | open |
| 3 | `analysis.go:94` — Ohm's Law check dead | open |
| 4 | `inference.go:113` — decode errors swallowed | open |
| 5 | `mappings.go:90` — `TempCoeffPPM` unused duplicate | open |
| 6 | `smd.go:255` — `findEIA96Multiplier` dead code | open |
| 7 | `bands.go:22` — wrong error message | open |
| 8 | `internal/cli/format.go:16` — `strings.Title` deprecated | open |
| 9 | `smd_test.go:135` — G115 integer overflow, nonsensical test names | open |
| 10 | `integration_test.go:27,42` — G204 false positives, need nolint | open |
| 11 | `cmd/resistor-cli/cmd/helpers.go:31` — De Morgan's simplification | open |
| 12 | `common.go:77,89` + `select.go:207` — WriteString+Sprintf pattern | open |
| 13 | TUI stubs — 9 unused symbols until Analyze/SMD views implemented | open |
