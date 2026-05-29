# syntax=docker/dockerfile:1
#
# Build args:
#   WASM=go      (default) standard Go toolchain — ~3.4 MB uncompressed
#   WASM=tinygo  opt-in    TinyGo toolchain       — ~1.1 MB uncompressed (~430 KB gzip)
#
# Usage:
#   docker build .                            # standard Go WASM
#   docker build --build-arg WASM=tinygo .   # TinyGo WASM
#
# The ARG must be declared before the first FROM so it can be used in FROM lines.
ARG WASM=go
ARG VERSION=dev

# ── Stage 1a: WASM via standard Go ───────────────────────────────────────────
FROM golang:1.26 AS wasm-go
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN GOOS=js GOARCH=wasm go build -o web/resistor.wasm ./cmd/resistor-wasm && \
    cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" web/wasm_exec.js

# ── Stage 1b: WASM via TinyGo ────────────────────────────────────────────────
FROM tinygo/tinygo:0.41.1 AS wasm-tinygo
USER root
WORKDIR /src
COPY go.mod go.sum ./
COPY . .
RUN tinygo build -target=wasm -o web/resistor.wasm ./cmd/resistor-wasm && \
    cp "$(tinygo env TINYGOROOT)/targets/wasm_exec.js" web/wasm_exec.js

# ── Stage 2: select WASM source based on build arg ───────────────────────────
FROM wasm-${WASM} AS wasm

# ── Stage 3: compile server binary ───────────────────────────────────────────
FROM golang:1.26 AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Overwrite the repo WASM artifacts with whichever variant was selected above.
COPY --from=wasm /src/web/resistor.wasm web/resistor.wasm
COPY --from=wasm /src/web/wasm_exec.js  web/wasm_exec.js
ARG VERSION=dev
RUN CGO_ENABLED=0 GOOS=linux go build \
      -ldflags="-s -w -X 'main.version=${VERSION}'" \
      -trimpath \
      -o /out/resistor-server \
      ./cmd/resistor-server

# ── Stage 4: minimal runtime (scratch) ───────────────────────────────────────
FROM scratch
# Run as an unprivileged UID (no /etc/passwd required in scratch).
USER 10001:10001
COPY --from=builder /out/resistor-server /resistor-server
EXPOSE 8080
ENTRYPOINT ["/resistor-server"]
