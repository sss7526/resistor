# Contributing

## Before You Start

Open an issue before writing code for anything non-trivial. It avoids duplicated effort and makes sure the change fits the project's direction before you invest time in it.

## Workflow

Fork the repository, create a branch for your change, and open a pull request against `main`.

```
git checkout -b your-branch-name
```

Branch naming is flexible. Use something that reflects what the branch does.

## CI Gates

Before opening a pull request, both of the following must pass:

```
make test-all
make smoke
```

`make test-all` runs unit tests and CLI integration tests. `make smoke` builds the CLI and runs a basic end-to-end check against the binary. A PR that fails either of these will not be reviewed.

## Code Style

Run `go fmt` before committing:

```
make fmt
```

The project has no linter configuration beyond standard formatting. Keep changes focused. A fix does not need surrounding cleanup, and a new feature does not need speculative abstractions.

## Commit Messages

Write commit messages in the imperative mood and describe why the change is being made, not just what it does. First line should be 72 characters or fewer.

Good:
```
fix parseESeries returning error on empty input
```

Less useful:
```
updated helpers.go
```

## Adding or Changing Public API

The core library (`github.com/sss7526/resistor`) is the stable, importable surface. Changes to exported types, function signatures, or behavior need to be considered carefully because they affect downstream consumers.

If you are adding a new exported symbol, make sure it fits the existing design:

- Deterministic operations must not perform inference.
- Inference operations must record all assumptions in the result.
- Confidence values must be in [0.0, 1.0].
- No silent defaults. Every default applied must be recorded in an `Assumptions` field.

## Tests

New deterministic logic requires unit tests. New inference rules require tests that verify the rule fires under the right conditions and does not fire under others. Fuzz targets exist for the primary decode paths; if you add a new decode function, add a fuzz target for it.

## Pull Requests

Keep pull requests focused on one thing. A PR that fixes a bug and refactors unrelated code is harder to review than two separate PRs.

Write a clear description of what the PR does and why. Reference any related issues.
