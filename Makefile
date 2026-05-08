# ---------------------------------------
# Go Resistor Library Makefile
# ---------------------------------------

# Package name (auto-detect)
PKG := ./...

# Default fuzz time (override via: make fuzz FUZZTIME=30s)
FUZZTIME ?= 10s

# ---------------------------------------
# Core Targets
# ---------------------------------------

.PHONY: all
all: fmt test

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

# ---------------------------------------
# Help
# ---------------------------------------

.PHONY: help
help:
	@echo ""
	@echo "Available targets:"
	@echo "  make fmt            - Run go fmt"
	@echo "  make test           - Run all unit tests"
	@echo "  make test-short     - Run short tests"
	@echo "  make test-bands     - Run band tests only"
	@echo "  make test-smd       - Run SMD tests only"
	@echo "  make test-series    - Run E-series tests only"
	@echo "  make test-select    - Run selection tests only"
	@echo "  make test-infer     - Run inference tests only"
	@echo "  make fuzz           - Run all fuzz tests"
	@echo "  make fuzz-bands     - Fuzz DecodeBands"
	@echo "  make fuzz-smd       - Fuzz DecodeSMD"
	@echo "  make fuzz-series    - Fuzz NearestStandard"
	@echo "  make clean          - Clean test cache & fuzz artifacts"
	@echo ""