package main

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var binary string

// TestMain builds the CLI binary once before running tests.
func TestMain(m *testing.M) {
	tmpDir, err := os.MkdirTemp("", "resistor-cli-test")
	if err != nil {
		panic(err)
	}

	binaryPath := filepath.Join(tmpDir, "resistor-cli")

	// Build the CLI from current directory (cmd/resistor-cli)
	ldflags := "-X github.com/sss7526/resistor/cmd/resistor-cli/cmd.version=v0.1.0"
	build := exec.Command("go", "build", "-ldflags", ldflags, "-o", binaryPath, ".") //nolint:gosec // intentional subprocess in test setup
	build.Dir = "." // current directory
	if err := build.Run(); err != nil {
		panic(err)
	}

	binary = binaryPath

	code := m.Run()

	_ = os.RemoveAll(tmpDir)
	os.Exit(code)
}

func runCLI(t *testing.T, args ...string) (string, error) {
	cmd := exec.Command(binary, args...) //nolint:gosec // intentional subprocess in test helper
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func runCLISuccess(t *testing.T, args ...string) string {
	out, err := runCLI(t, args...)
	require.NoError(t, err, out)
	return out
}

func runCLIError(t *testing.T, args ...string) string {
	out, err := runCLI(t, args...)
	require.Error(t, err)
	return out
}

///////////////////////////////////////////////////////////////////////////////
// Version
///////////////////////////////////////////////////////////////////////////////

func TestCLI_Version(t *testing.T) {
	out := runCLISuccess(t, "version")
	require.Contains(t, out, "v0.1.0")
}

///////////////////////////////////////////////////////////////////////////////
// Select Command
///////////////////////////////////////////////////////////////////////////////

func TestCLI_Select_Default(t *testing.T) {
	out := runCLISuccess(t, "select", "487")
	require.Contains(t, out, "Selected")
	require.Contains(t, out, "Bands")
}

func TestCLI_Select_WithSeries(t *testing.T) {
	out := runCLISuccess(t, "select", "487", "--series", "E12")
	require.Contains(t, out, "Selected")
}

func TestCLI_Select_JSON(t *testing.T) {
	out := runCLISuccess(t, "select", "487", "--json")

	var parsed map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(out), &parsed))
	require.Equal(t, true, parsed["success"])
}

func TestCLI_Select_InvalidSeries(t *testing.T) {
	out := runCLIError(t, "select", "487", "--series", "E999")
	require.Contains(t, out, "invalid E-series")
}

///////////////////////////////////////////////////////////////////////////////
// Infer Command
///////////////////////////////////////////////////////////////////////////////

func TestCLI_Infer_Bands(t *testing.T) {
	out := runCLISuccess(t, "infer", "--bands", "brown,black,red,gold")
	require.Contains(t, out, "Resistance")
	require.Contains(t, out, "Tolerance")
}

func TestCLI_Infer_SMD(t *testing.T) {
	out := runCLISuccess(t, "infer", "--smd", "472")
	require.Contains(t, out, "4700")
}

func TestCLI_Infer_InvalidBand(t *testing.T) {
	out := runCLIError(t, "infer", "--bands", "brown,banana,red")
	require.Contains(t, out, "invalid band color")
}

func TestCLI_Infer_JSON(t *testing.T) {
	out := runCLISuccess(t, "infer", "--bands", "brown,black,red,gold", "--json")

	var parsed map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(out), &parsed))
	require.Equal(t, true, parsed["success"])
}

///////////////////////////////////////////////////////////////////////////////
// Analyze Command
///////////////////////////////////////////////////////////////////////////////

func TestCLI_Analyze_VoltageDriven(t *testing.T) {
	out := runCLISuccess(t, "analyze", "--r", "100", "--v", "10", "--pwr", "0.5")
	require.Contains(t, out, "Voltage")
	require.Contains(t, out, "Power")
	require.Contains(t, out, "Derated Safe")
}

func TestCLI_Analyze_WorstCaseBounds(t *testing.T) {
	out := runCLISuccess(t, "analyze", "--r", "100", "--v", "10", "--tol", "5")
	require.Contains(t, out, "R Min (WC)")
	require.Contains(t, out, "R Max (WC)")
}

func TestCLI_Analyze_NoDeratedSafe_WithoutPwr(t *testing.T) {
	out := runCLISuccess(t, "analyze", "--r", "100", "--v", "10")
	require.NotContains(t, out, "Derated Safe")
}

func TestCLI_Analyze_NoWorstCase_WithoutTol(t *testing.T) {
	out := runCLISuccess(t, "analyze", "--r", "100", "--v", "10")
	require.NotContains(t, out, "R Min (WC)")
	require.NotContains(t, out, "R Max (WC)")
}

func TestCLI_Analyze_NegativeVoltage(t *testing.T) {
	out := runCLIError(t, "analyze", "--r", "100", "--v", "-10")
	require.Contains(t, out, "voltage")
}

func TestCLI_Analyze_NegativeCurrent(t *testing.T) {
	out := runCLIError(t, "analyze", "--r", "100", "--i", "-0.1")
	require.Contains(t, out, "current")
}

func TestCLI_Analyze_NegativePower(t *testing.T) {
	out := runCLIError(t, "analyze", "--r", "100", "--v", "10", "--pwr", "-0.5")
	require.Contains(t, out, "power")
}

func TestCLI_Analyze_NegativeTolerance(t *testing.T) {
	out := runCLIError(t, "analyze", "--r", "100", "--v", "10", "--tol", "-5")
	require.Contains(t, out, "tolerance")
}

func TestCLI_Analyze_ToleranceAbove100(t *testing.T) {
	out := runCLIError(t, "analyze", "--r", "100", "--v", "10", "--tol", "150")
	require.Contains(t, out, "tolerance")
}

func TestCLI_Analyze_FullTolerance(t *testing.T) {
	// 100% tolerance: R min = 0, R max = 2R — both bounds must appear
	out := runCLISuccess(t, "analyze", "--r", "100", "--v", "10", "--tol", "100")
	require.Contains(t, out, "R Min (WC)")
	require.Contains(t, out, "R Max (WC)")
}

func TestCLI_Analyze_CurrentDriven(t *testing.T) {
	out := runCLISuccess(t, "analyze", "--r", "50", "--i", "0.2")
	require.Contains(t, out, "Voltage")
}

func TestCLI_Analyze_Invalid(t *testing.T) {
	out := runCLIError(t, "analyze", "--v", "10")
	require.Contains(t, out, "resistance")
}

func TestCLI_Analyze_JSON(t *testing.T) {
	out := runCLISuccess(t, "analyze", "--r", "100", "--v", "10", "--json")

	var parsed map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(out), &parsed))
	require.Equal(t, true, parsed["success"])
}

func TestCLI_Analyze_JSON_OptionalFieldsPresent(t *testing.T) {
	out := runCLISuccess(t, "analyze", "--r", "100", "--v", "10", "--pwr", "0.5", "--tol", "5", "--json")

	var parsed map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(out), &parsed))

	data, ok := parsed["data"].(map[string]interface{})
	require.True(t, ok)

	require.InDelta(t, 0.25, data["DeratedSafePower"], 1e-9)
	require.InDelta(t, 95.0, data["WorstCaseResistanceMin"], 1e-9)
	require.InDelta(t, 105.0, data["WorstCaseResistanceMax"], 1e-9)
}

func TestCLI_Analyze_JSON_OptionalFieldsAbsent(t *testing.T) {
	out := runCLISuccess(t, "analyze", "--r", "100", "--v", "10", "--json")

	var parsed map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(out), &parsed))

	data, ok := parsed["data"].(map[string]interface{})
	require.True(t, ok)

	_, hasDerated := data["DeratedSafePower"]
	require.False(t, hasDerated, "DeratedSafePower should be absent from JSON when --pwr is not provided")

	_, hasWCMin := data["WorstCaseResistanceMin"]
	require.False(t, hasWCMin, "WorstCaseResistanceMin should be absent from JSON when --tol is not provided")

	_, hasWCMax := data["WorstCaseResistanceMax"]
	require.False(t, hasWCMax, "WorstCaseResistanceMax should be absent from JSON when --tol is not provided")
}

///////////////////////////////////////////////////////////////////////////////
// SMD Command
///////////////////////////////////////////////////////////////////////////////

func TestCLI_SMD_Decode(t *testing.T) {
	out := runCLISuccess(t, "smd", "decode", "472")
	require.Contains(t, out, "4700")
}

func TestCLI_SMD_Encode(t *testing.T) {
	out := runCLISuccess(t, "smd", "encode", "4700")
	require.Contains(t, out, "472")
}

func TestCLI_SMD_Invalid(t *testing.T) {
	out := runCLIError(t, "smd", "decode", "XYZ")
	require.True(t, strings.Contains(out, "unsupported") || strings.Contains(out, "invalid"))
}

func TestCLI_SMD_JSON(t *testing.T) {
	out := runCLISuccess(t, "smd", "decode", "472", "--json")

	var parsed map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(out), &parsed))
	require.Equal(t, true, parsed["success"])
}

///////////////////////////////////////////////////////////////////////////////
// JSON exit code on error
///////////////////////////////////////////////////////////////////////////////

func TestCLI_JSON_ExitCodeOnError(t *testing.T) {
	cmd := exec.Command(binary, "smd", "decode", "XYZ", "--json") //nolint:gosec // intentional subprocess in test
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	err := cmd.Run()

	require.Error(t, err, "command should exit non-zero on failure with --json")

	exitErr, ok := err.(*exec.ExitError)
	require.True(t, ok, "expected *exec.ExitError")
	require.NotEqual(t, 0, exitErr.ExitCode(), "exit code should be non-zero")

	var parsed map[string]interface{}
	require.NoError(t, json.Unmarshal(stdout.Bytes(), &parsed))
	require.Equal(t, false, parsed["success"])
	require.NotEmpty(t, parsed["error"])
}
