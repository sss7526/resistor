# Changelog

## [1.0.4](https://github.com/sss7526/resistor/compare/v1.0.3...v1.0.4) (2026-05-29)


### Bug Fixes

* address all 8 confirmed code-review findings ([da94ef0](https://github.com/sss7526/resistor/commit/da94ef00d8ff52527fd790886deac9a5c3a27c8c))

## [1.0.3](https://github.com/sss7526/resistor/compare/v1.0.2...v1.0.3) (2026-05-29)


### Bug Fixes

* set DOCKER_API_VERSION for watchtower compatibility ([bfb03ea](https://github.com/sss7526/resistor/commit/bfb03ea9ab28ab0e3f115c86a7fccde1071147ef))
* set explicit container_name for watchtower to match ([8288c5b](https://github.com/sss7526/resistor/commit/8288c5bf57cec9230b210b952e65840afc622cf8))

## [1.0.2](https://github.com/sss7526/resistor/compare/v1.0.1...v1.0.2) (2026-05-29)


### Bug Fixes

* commit empty WASM stubs so go vet resolves embed patterns ([9ab491e](https://github.com/sss7526/resistor/commit/9ab491e487942ebd12dbd50e7086bf2a8ddcc6c4))
* correct misspelling of "unknown" in PackageType constant ([4cf49d9](https://github.com/sss7526/resistor/commit/4cf49d97def773fc6e8e2febb27bc41f1f98c5bb))
* make lint depend on build-wasm ([adda309](https://github.com/sss7526/resistor/commit/adda30983a3242cc2bc9539afdbc0fe3cac1570c))
* match embed lint exclusion on error text not file path ([48cc056](https://github.com/sss7526/resistor/commit/48cc05603ac482742a4e2a18cba72e542e5a10dd))
* update golangci.yml to v2 config format ([e0564d7](https://github.com/sss7526/resistor/commit/e0564d79795dacce4c9197f86ad73b978f8adb25))

## [1.0.1](https://github.com/sss7526/resistor/compare/v1.0.0...v1.0.1) (2026-05-29)


### Bug Fixes

* run wasm-tinygo build stage as root ([4eea790](https://github.com/sss7526/resistor/commit/4eea790764e5ccc4175f6d4eb7653aa3c2ee132f))

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
