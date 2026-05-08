// Package resistor provides deterministic standards support and
// heuristic inference capabilities for fixed resistors.
//
// The package is organized into layered responsibilities:
//
//  1. Standards Encoding/Decoding
//     - 4, 5, and 6-band color code support (IEC 60062)
//     - SMD 3-digit, 4-digit, R-notation, and EIA‑96 support
//     - Full E-series (E3–E192) preferred value support (IEC 60063)
//
//  2. Selection Engine
//     - Snap arbitrary resistance values to preferred E-series values
//     - Encode selected values into visual band representations
//     - Explicitly track defaults and assumptions
//
//  3. Inference Engine
//     - Deterministic extraction of known properties
//     - Heuristic inference of type, power rating, voltage rating
//     - Monotonic confidence model
//     - Explicit assumption tracking
//
// Design Principles:
//
//   - Deterministic logic always overrides heuristic inference.
//   - All assumptions are explicitly recorded.
//   - Confidence values are bounded in [0.0, 1.0].
//   - No silent guessing or hidden defaults.
//   - API symmetry between encoding and decoding where possible.
//
// This package is suitable for:
//   - Engineering tools
//   - Educational applications
//   - Interactive TUI interfaces
//   - WASM/web frontends
//
// The API surface is intentionally explicit and structured to allow
// long-term stability and future extensibility.
package resistor
