# ---------------------------------------
# Go Resistor Library Makefile
# ---------------------------------------

GOBIN := $(shell go env GOBIN)
ifeq ($(GOBIN),)
GOBIN := $(shell go env GOPATH)/bin
endif

# Package name (auto-detect)
PKG := ./...
BINDIR := bin
CLI := resistor-cli
CLI_PATH := ./cmd/resistor-cli
TUI := resistor-tui
TUI_PATH := ./cmd/resistor-tui
WASM_PATH := ./cmd/resistor-wasm
SERVER := resistor-server
SERVER_PATH := ./cmd/resistor-server
WEBDIR := web
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

# TinyGo WASM build (pass TINYGO=1 or use build-wasm-tinygo)
TINYGO ?= 0
TINYGO_BIN ?= tinygo

# Default fuzz time (override via: make fuzz FUZZTIME=30s)
FUZZTIME ?= 10s

# ---------------------------------------
# Core Targets
# ---------------------------------------

.PHONY: all
all: fmt test build

.PHONY: build
build: build-cli build-tui build-wasm build-server

# Format code
.PHONY: fmt
fmt:
	@echo "→ Running go fmt"
	go fmt $(PKG)

# Run all unit tests
.PHONY: test
test:
	@echo "→ Running unit tests"
	go test -v $(PKG)

# Run short unit tests
.PHONY: test-short
test-short:
	@echo "→ Running short unit tests"
	go test -short -v $(PKG)

.PHONY: test-cli
test-cli: build
	go test ./cmd/resistor-cli -v

.PHONY: test-all
test-all: test test-cli

.PHONY: smoke
smoke: build
	./$(BINDIR)/$(CLI) select 487 > /dev/null
	./$(BINDIR)/$(CLI) infer --bands brown,black,red,gold > /dev/null
	./$(BINDIR)/$(CLI) analyze --r 100 --v 10 > /dev/null
	./$(BINDIR)/$(CLI) smd decode 472 > /dev/null
	@echo "Smoke tests passed"

# ---------------------------------------
# CLI Build Targets
# ---------------------------------------

.PHONY: build-cli
build-cli:
	@echo "→ Building CLI binary"
	@mkdir -p $(BINDIR)
	go build -o $(BINDIR)/$(CLI) -ldflags "-X 'github.com/sss7526/resistor/cmd/resistor-cli/cmd.version=$(VERSION)'" $(CLI_PATH)

.PHONY: install-cli
install-cli:
	@echo "→ Installing CLI"
	go install $(CLI_PATH)

# ---------------------------------------
# TUI Build Targets
# ---------------------------------------

.PHONY: build-tui
build-tui:
	@echo "→ Building TUI binary"
	@mkdir -p $(BINDIR)
	go build -o $(BINDIR)/$(TUI) -ldflags "-X 'github.com/sss7526/resistor/cmd/resistor-tui/app.version=$(VERSION)'" $(TUI_PATH)

.PHONY: install-tui
install-tui:
	@echo "→ Installing TUI"
	go install $(TUI_PATH)

# ---------------------------------------
# Server Build Targets
# ---------------------------------------

.PHONY: build-server
build-server: build-wasm
	@echo "→ Building server binary"
	@mkdir -p $(BINDIR)
	go build -o $(BINDIR)/$(SERVER) -ldflags "-X 'main.version=$(VERSION)'" $(SERVER_PATH)

# Convenience alias: build server with TinyGo WASM (~430 KB gzip vs ~3.4 MB)
.PHONY: build-server-tinygo
build-server-tinygo:
	$(MAKE) build-server TINYGO=1

.PHONY: install-server
install-server:
	@echo "→ Installing server"
	go install $(SERVER_PATH)

.PHONY: install-server-tinygo
install-server-tinygo:
	$(MAKE) build-wasm TINYGO=1
	@echo "→ Installing server (TinyGo WASM)"
	go install $(SERVER_PATH)

# ---------------------------------------
# WASM Build Targets
# ---------------------------------------

.PHONY: build-wasm
ifeq ($(TINYGO),1)
build-wasm:
	@echo "→ Building WASM module (TinyGo)"
	@mkdir -p $(WEBDIR)
	$(TINYGO_BIN) build -target=wasm -o $(WEBDIR)/resistor.wasm $(WASM_PATH)
	cp "$$($(TINYGO_BIN) env TINYGOROOT)/targets/wasm_exec.js" $(WEBDIR)/wasm_exec.js
	@echo "→ Artifacts: $(WEBDIR)/resistor.wasm  $(WEBDIR)/wasm_exec.js"
else
build-wasm:
	@echo "→ Building WASM module (standard Go)"
	@mkdir -p $(WEBDIR)
	GOOS=js GOARCH=wasm go build -o $(WEBDIR)/resistor.wasm $(WASM_PATH)
	cp "$$(go env GOROOT)/lib/wasm/wasm_exec.js" $(WEBDIR)/wasm_exec.js
	@echo "→ Artifacts: $(WEBDIR)/resistor.wasm  $(WEBDIR)/wasm_exec.js"
endif

.PHONY: build-wasm-tinygo
build-wasm-tinygo:
	$(MAKE) build-wasm TINYGO=1

# ---------------------------------------
# Docker Build Targets
# ---------------------------------------

.PHONY: docker-build
docker-build:
	@echo "→ Building Docker image (standard Go WASM)"
	docker build --build-arg WASM=go --build-arg VERSION=$(VERSION) \
	  -t resistor-server:$(VERSION) .

.PHONY: docker-build-tinygo
docker-build-tinygo:
	@echo "→ Building Docker image (TinyGo WASM)"
	docker build --build-arg WASM=tinygo --build-arg VERSION=$(VERSION) \
	  -t resistor-server:$(VERSION)-tinygo .

# ---------------------------------------
# Unit Test Subsets
# ---------------------------------------

.PHONY: test-bands
test-bands:
	go test -v -run TestDecodeBands

.PHONY: test-smd
test-smd:
	go test -v -run TestSMD

.PHONY: test-series
test-series:
	go test -v -run TestNearestStandard

.PHONY: test-select
test-select:
	go test -v -run TestSelectStandardResistor

.PHONY: test-infer
test-infer:
	go test -v -run TestInferResistor

.PHONY: test-analyze
test-analyze:
	go test -v -run TestAnalyzeResistor

# ---------------------------------------
# Fuzz Targets
# ---------------------------------------

.PHONY: fuzz-bands
fuzz-bands:
	@echo "→ Fuzzing DecodeBands"
	go test -fuzz=FuzzDecodeBands -fuzztime=$(FUZZTIME)

.PHONY: fuzz-smd
fuzz-smd:
	@echo "→ Fuzzing DecodeSMD"
	go test -fuzz=FuzzDecodeSMD -fuzztime=$(FUZZTIME)

.PHONY: fuzz-series
fuzz-series:
	@echo "→ Fuzzing NearestStandard"
	go test -fuzz=FuzzNearestStandard -fuzztime=$(FUZZTIME)

.PHONY: fuzz
fuzz: fuzz-bands fuzz-smd fuzz-series

# ---------------------------------------
# Lint
# ---------------------------------------

.PHONY: lint
lint:
	@echo "→ Running golangci-lint"
	golangci-lint run --verbose --color always --enable-only=errcheck,govet,ineffassign,staticcheck,unused,gosec --timeout=5m

.PHONY: vuln
vuln:
	@echo "→ Running govulncheck"
	govulncheck -show verbose,color ./...

# ---------------------------------------
# Clean
# ---------------------------------------

.PHONY: clean
clean:
	@echo "→ Cleaning fuzz cache"
	go clean -testcache
	rm -rf testdata/fuzz
	rm -rf $(BINDIR)
	rm -f $(GOBIN)/$(CLI)
	rm -f $(GOBIN)/$(TUI)
	rm -f $(GOBIN)/$(SERVER)
	rm -f $(WEBDIR)/resistor.wasm $(WEBDIR)/wasm_exec.js

# ---------------------------------------
# Help
# ---------------------------------------

.PHONY: help
help:
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Default target: all (fmt + test + build)"
	@echo ""
	@echo "Build"
	@echo "  build               Build CLI, TUI, WASM, and server"
	@echo "  build-cli           Build CLI binary to bin/"
	@echo "  build-tui           Build TUI binary to bin/"
	@echo "  build-wasm          Build WASM module to web/ (standard Go)"
	@echo "  build-wasm TINYGO=1 Build WASM module using TinyGo (~430 KB gzip)"
	@echo "  build-wasm-tinygo   Alias for build-wasm TINYGO=1"
	@echo "  build-server          Build server (standard Go WASM)"
	@echo "  build-server TINYGO=1 Build server with TinyGo WASM (~430 KB gzip)"
	@echo "  build-server-tinygo   Alias for build-server TINYGO=1"
	@echo "  install-cli           Install CLI to GOPATH/bin"
	@echo "  install-tui           Install TUI to GOPATH/bin"
	@echo "  install-server        Install server to GOPATH/bin (standard Go WASM)"
	@echo "  install-server-tinygo Install server to GOPATH/bin (TinyGo WASM)"
	@echo "  docker-build          Build Docker image with standard Go WASM"
	@echo "  docker-build-tinygo   Build Docker image with TinyGo WASM (:latest tag)"
	@echo ""
	@echo "Test"
	@echo "  test                Unit tests"
	@echo "  test-short          Unit tests (short mode)"
	@echo "  test-cli            CLI integration tests (builds first)"
	@echo "  test-all            Unit tests + CLI integration tests"
	@echo "  smoke               End-to-end smoke tests against built binary"
	@echo ""
	@echo "Test subsets"
	@echo "  test-bands          Band encode/decode tests"
	@echo "  test-smd            SMD encode/decode tests"
	@echo "  test-series         E-series tests"
	@echo "  test-select         Selection engine tests"
	@echo "  test-infer          Inference engine tests"
	@echo "  test-analyze        Analysis engine tests"
	@echo ""
	@echo "Fuzz  (override duration: make fuzz FUZZTIME=60s, default: $(FUZZTIME))"
	@echo "  fuzz                Run all fuzz targets"
	@echo "  fuzz-bands          Fuzz DecodeBands"
	@echo "  fuzz-smd            Fuzz DecodeSMD"
	@echo "  fuzz-series         Fuzz NearestStandard"
	@echo ""
	@echo "Other"
	@echo "  fmt                 Run go fmt"
	@echo "  lint                Run golangci-lint (errcheck,govet,ineffassign,staticcheck,unused,gosec)"
	@echo "  clean               Remove bin/, test cache, and fuzz artifacts"
	@echo "  vuln                Run govulncheck"
	@echo "  help                Show this message"
	@echo ""