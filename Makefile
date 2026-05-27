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
VERSION := v0.1.0

# Default fuzz time (override via: make fuzz FUZZTIME=30s)
FUZZTIME ?= 10s

# ---------------------------------------
# Core Targets
# ---------------------------------------

.PHONY: all
all: fmt test build

.PHONY: build
build: build-cli build-tui

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

# ---------------------------------------
# Help
# ---------------------------------------

.PHONY: help
help:
	@echo ""
	@echo "Available targets:"
	@echo "  make fmt            - Run go fmt"
	@echo "  make build          - Build both binaries to bin/"
	@echo "  make build-cli      - Build CLI binary to bin/"
	@echo "  make build-tui      - Build TUI binary to bin/"
	@echo "  make install-cli    - Install CLI to GOPATH/bin"
	@echo "  make install-tui    - Install TUI to GOPATH/bin"
	@echo "  make test           - Run all unit tests"
	@echo "  make test-short     - Run short tests"
	@echo "  make test-bands     - Run band tests only"
	@echo "  make test-smd       - Run SMD tests only"
	@echo "  make test-series    - Run E-series tests only"
	@echo "  make test-select    - Run selection tests only"
	@echo "  make test-infer     - Run inference tests only"
	@echo "  make test-analyze   - Run analysis tests only"
	@echo "  make fuzz           - Run all fuzz tests"
	@echo "  make fuzz-bands     - Fuzz DecodeBands"
	@echo "  make fuzz-smd       - Fuzz DecodeSMD"
	@echo "  make fuzz-series    - Fuzz NearestStandard"
	@echo "  make clean          - Clean binary & test cache & fuzz artifacts"
	@echo ""