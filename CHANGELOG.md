# Changelog

## 1.0.0 (2026-05-28)


### Features

* add SVG favicon (10 kΩ resistor, dark theme) ([965e093](https://github.com/sss7526/resistor/commit/965e093983a94d3199eaaee73918f1e24a7da8bb))
* CTA footer (M15), containerised deployment (M16) ([8124797](https://github.com/sss7526/resistor/commit/8124797cf89562632918facc1fb6248ec10e5811))
* implement M17 — CI, release-please, GoReleaser, Dependabot, version wiring ([8b58a4d](https://github.com/sss7526/resistor/commit/8b58a4dfb2c99d3eab6599a7c76306d978efb7a8))
* implement M17 automation (CI, release-please, GoReleaser, Dependabot) ([bb96384](https://github.com/sss7526/resistor/commit/bb963849ea18937c96fb921a6243c21bbffbd1da))
* replace resistor favicon with Ω symbol ([f3bee94](https://github.com/sss7526/resistor/commit/f3bee943d0d22b6135790dce400eec38a8474dea))


### Bug Fixes

* add bump-minor-pre-major to release-please config ([dc40752](https://github.com/sss7526/resistor/commit/dc407529f88fae77bbd69bf071e6cc536eb51133))
* add wasm-unsafe-eval to CSP script-src for WASM instantiation ([b5e6281](https://github.com/sss7526/resistor/commit/b5e628121beb5e9d9b7dda4ddbf4f628d110335a))
* build-wasm before running tests ([c639f18](https://github.com/sss7526/resistor/commit/c639f184cf60a3403c4fdf231d81f34eafd80253))
* compose file pulls from Docker Hub instead of building from source ([ffbe269](https://github.com/sss7526/resistor/commit/ffbe269cc9271c9afdb4a260ab94ab2a0903dfde))
* remove invalid skip-github-release input; fix goreleaser docker ids ([bd7743a](https://github.com/sss7526/resistor/commit/bd7743a49c65f21c035827f39739bc6ccd2bb54e))
* replace inline onclick handlers with addEventListener (CSP script-src-attr) ([64a0710](https://github.com/sss7526/resistor/commit/64a0710de8d6675e40091b9ba0d322bc5e01b1fc))
