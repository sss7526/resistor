//go:build js && wasm

package main

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/sss7526/resistor"
)

func main() {
	obj := js.Global().Get("Object").New()

	obj.Set("decodeBands", js.FuncOf(safe(decodeBandsJS)))
	obj.Set("encodeBands", js.FuncOf(safe(encodeBandsJS)))
	obj.Set("decodeSMD", js.FuncOf(safe(decodeSMDJS)))
	obj.Set("encodeSMD", js.FuncOf(safe(encodeSMDJS)))
	obj.Set("nearestStandard", js.FuncOf(safe(nearestStandardJS)))
	obj.Set("selectStandardResistor", js.FuncOf(safe(selectStandardResistorJS)))
	obj.Set("inferResistor", js.FuncOf(safe(inferResistorJS)))
	obj.Set("analyzeResistor", js.FuncOf(safe(analyzeResistorJS)))

	js.Global().Set("resistor", obj)

	<-make(chan struct{})
}

///////////////////////////////////////////////////////////////////////////////
// Response helpers
///////////////////////////////////////////////////////////////////////////////

func jsOK(v any) js.Value {
	type resp struct {
		OK    bool `json:"ok"`
		Value any  `json:"value"`
	}
	b, _ := json.Marshal(resp{OK: true, Value: v})
	return js.Global().Get("JSON").Call("parse", string(b))
}

func jsErr(err error) js.Value {
	type resp struct {
		OK    bool   `json:"ok"`
		Error string `json:"error"`
	}
	b, _ := json.Marshal(resp{OK: false, Error: err.Error()})
	return js.Global().Get("JSON").Call("parse", string(b))
}

// safe wraps a handler so panics surface as JS errors instead of crashing.
func safe(fn func(js.Value, []js.Value) js.Value) func(js.Value, []js.Value) any {
	return func(this js.Value, args []js.Value) (ret any) {
		defer func() {
			if r := recover(); r != nil {
				ret = jsErr(fmt.Errorf("internal error: %v", r))
			}
		}()
		return fn(this, args)
	}
}

// argString returns args[i].String() or "" if out of range.
func argString(args []js.Value, i int) string {
	if i >= len(args) || args[i].IsNull() || args[i].IsUndefined() {
		return ""
	}
	return args[i].String()
}

// argJSON returns the first argument as a raw JSON string.
func argJSON(args []js.Value) (string, bool) {
	if len(args) == 0 || args[0].IsNull() || args[0].IsUndefined() {
		return "", false
	}
	return args[0].String(), true
}

///////////////////////////////////////////////////////////////////////////////
// decodeBands(jsonInput: string) → {ok, value: ResistorSpec}
//
// jsonInput: JSON array of color strings
// Example:  decodeBands('["green","brown","brown","gold"]')
///////////////////////////////////////////////////////////////////////////////

func decodeBandsJS(_ js.Value, args []js.Value) js.Value {
	raw, ok := argJSON(args)
	if !ok {
		return jsErr(fmt.Errorf("expected JSON array of color strings"))
	}
	var colorNames []string
	if err := json.Unmarshal([]byte(raw), &colorNames); err != nil {
		return jsErr(fmt.Errorf("invalid input: %w", err))
	}
	bands := make([]resistor.Color, len(colorNames))
	for i, n := range colorNames {
		bands[i] = resistor.Color(n)
	}
	spec, err := resistor.DecodeBands(bands)
	if err != nil {
		return jsErr(err)
	}
	return jsOK(spec)
}

///////////////////////////////////////////////////////////////////////////////
// encodeBands(jsonInput: string) → {ok, value: string[]}
//
// jsonInput: JSON object with resistor spec fields
// Example:  encodeBands('{"resistanceOhms":510,"tolerancePct":5}')
///////////////////////////////////////////////////////////////////////////////

type jsSpecInput struct {
	ResistanceOhms float64 `json:"resistanceOhms"`
	TolerancePct   float64 `json:"tolerancePct"`
	PowerWatts     float64 `json:"powerWatts"`
	TempCoeffPPM   int     `json:"tempCoeffPPM"`
	Package        string  `json:"package"`
	Type           string  `json:"type"`
}

func encodeBandsJS(_ js.Value, args []js.Value) js.Value {
	raw, ok := argJSON(args)
	if !ok {
		return jsErr(fmt.Errorf("expected JSON spec object"))
	}
	var in jsSpecInput
	if err := json.Unmarshal([]byte(raw), &in); err != nil {
		return jsErr(fmt.Errorf("invalid input: %w", err))
	}
	pkg, _ := resistor.ParsePackageType(in.Package)
	spec := resistor.ResistorSpec{
		ResistanceOhms: in.ResistanceOhms,
		TolerancePct:   in.TolerancePct,
		PowerWatts:     in.PowerWatts,
		TempCoeffPPM:   in.TempCoeffPPM,
		Package:        pkg,
		Type:           resistor.ResistorType(in.Type),
	}
	bands, err := resistor.EncodeBands(spec)
	if err != nil {
		return jsErr(err)
	}
	names := make([]string, len(bands))
	for i, c := range bands {
		names[i] = string(c)
	}
	return jsOK(names)
}

///////////////////////////////////////////////////////////////////////////////
// decodeSMD(marking: string) → {ok, value: ResistorSpec}
//
// Example:  decodeSMD("472")
///////////////////////////////////////////////////////////////////////////////

func decodeSMDJS(_ js.Value, args []js.Value) js.Value {
	marking := argString(args, 0)
	if marking == "" {
		return jsErr(fmt.Errorf("marking string required"))
	}
	spec, err := resistor.DecodeSMD(marking)
	if err != nil {
		return jsErr(err)
	}
	return jsOK(spec)
}

///////////////////////////////////////////////////////////////////////////////
// encodeSMD(jsonInput: string) → {ok, value: string}
//
// jsonInput: {"resistance": number, "mode"?: "auto"|"standard"|"eia96"}
// Example:  encodeSMD('{"resistance":4700}')
///////////////////////////////////////////////////////////////////////////////

func encodeSMDJS(_ js.Value, args []js.Value) js.Value {
	raw, ok := argJSON(args)
	if !ok {
		return jsErr(fmt.Errorf("expected JSON input object"))
	}
	var in struct {
		Resistance float64 `json:"resistance"`
		Mode       string  `json:"mode"`
	}
	if err := json.Unmarshal([]byte(raw), &in); err != nil {
		return jsErr(fmt.Errorf("invalid input: %w", err))
	}
	mode, err := parseSMDMode(in.Mode)
	if err != nil {
		return jsErr(err)
	}
	marking, err := resistor.EncodeSMD(in.Resistance, mode)
	if err != nil {
		return jsErr(err)
	}
	return jsOK(marking)
}

func parseSMDMode(s string) (resistor.SMDEncodingMode, error) {
	switch s {
	case "", "auto":
		return resistor.SMDAuto, nil
	case "standard":
		return resistor.SMDStandard, nil
	case "eia96":
		return resistor.SMDEIA96, nil
	default:
		return 0, fmt.Errorf("unknown SMD mode %q: use auto, standard, or eia96", s)
	}
}

///////////////////////////////////////////////////////////////////////////////
// nearestStandard(jsonInput: string) → {ok, value: number}
//
// jsonInput: {"value": number, "series"?: string, "mode"?: string}
// Example:  nearestStandard('{"value":487,"series":"E24","mode":"nearest"}')
///////////////////////////////////////////////////////////////////////////////

func nearestStandardJS(_ js.Value, args []js.Value) js.Value {
	raw, ok := argJSON(args)
	if !ok {
		return jsErr(fmt.Errorf("expected JSON input object"))
	}
	var in struct {
		Value  float64 `json:"value"`
		Series string  `json:"series"`
		Mode   string  `json:"mode"`
	}
	if err := json.Unmarshal([]byte(raw), &in); err != nil {
		return jsErr(fmt.Errorf("invalid input: %w", err))
	}
	series, err := resistor.ParseESeries(in.Series)
	if err != nil {
		series = resistor.E24
	}
	mode, err := resistor.ParseRoundingMode(in.Mode)
	if err != nil {
		return jsErr(err)
	}
	result, err := resistor.NearestStandard(in.Value, series, mode)
	if err != nil {
		return jsErr(err)
	}
	return jsOK(result)
}

///////////////////////////////////////////////////////////////////////////////
// selectStandardResistor(jsonInput: string) → {ok, value: SelectionResult}
//
// jsonInput: {"resistance": number, "series"?: string, "tolerancePct"?: number, "rounding"?: string}
// Example:  selectStandardResistor('{"resistance":487,"series":"E24"}')
///////////////////////////////////////////////////////////////////////////////

func selectStandardResistorJS(_ js.Value, args []js.Value) js.Value {
	raw, ok := argJSON(args)
	if !ok {
		return jsErr(fmt.Errorf("expected JSON input object"))
	}
	var in struct {
		Resistance   float64 `json:"resistance"`
		Series       string  `json:"series"`
		TolerancePct float64 `json:"tolerancePct"`
		Rounding     string  `json:"rounding"`
	}
	if err := json.Unmarshal([]byte(raw), &in); err != nil {
		return jsErr(fmt.Errorf("invalid input: %w", err))
	}
	series, _ := resistor.ParseESeries(in.Series)
	rounding, _ := resistor.ParseRoundingMode(in.Rounding)
	req := resistor.SelectionRequest{
		Resistance:   in.Resistance,
		Series:       series,
		TolerancePct: in.TolerancePct,
		Rounding:     rounding,
	}
	result, err := resistor.SelectStandardResistor(req)
	if err != nil {
		return jsErr(err)
	}
	return jsOK(result)
}

///////////////////////////////////////////////////////////////////////////////
// inferResistor(jsonInput: string) → {ok, value: InferenceResult}
//
// jsonInput: {
//   "bands"?: string[],
//   "bodyColor"?: string,
//   "lengthMM"?: number,
//   "package"?: string,
//   "marking"?: string
// }
///////////////////////////////////////////////////////////////////////////////

func inferResistorJS(_ js.Value, args []js.Value) js.Value {
	raw, ok := argJSON(args)
	if !ok {
		return jsErr(fmt.Errorf("expected JSON input object"))
	}
	var in struct {
		Bands     []string `json:"bands"`
		BodyColor string   `json:"bodyColor"`
		LengthMM  float64  `json:"lengthMM"`
		Package   string   `json:"package"`
		Marking   string   `json:"marking"`
	}
	if err := json.Unmarshal([]byte(raw), &in); err != nil {
		return jsErr(fmt.Errorf("invalid input: %w", err))
	}
	pkg, _ := resistor.ParsePackageType(in.Package)
	obs := resistor.ObservedResistor{
		BodyColor: resistor.Color(in.BodyColor),
		LengthMM:  in.LengthMM,
		Package:   pkg,
		Marking:   in.Marking,
	}
	if len(in.Bands) > 0 {
		obs.Bands = make([]resistor.Color, len(in.Bands))
		for i, b := range in.Bands {
			obs.Bands[i] = resistor.Color(b)
		}
	}
	result, err := resistor.InferResistor(obs)
	if err != nil {
		return jsErr(err)
	}
	return jsOK(result)
}

///////////////////////////////////////////////////////////////////////////////
// analyzeResistor(jsonInput: string) → {ok, value: AnalysisReport}
//
// jsonInput: {
//   "spec": {
//     "resistanceOhms": number,
//     "powerWatts"?: number,
//     "tolerancePct"?: number,
//     "tempCoeffPPM"?: number
//   },
//   "appliedVoltage"?: number,
//   "appliedCurrent"?: number
// }
///////////////////////////////////////////////////////////////////////////////

func analyzeResistorJS(_ js.Value, args []js.Value) js.Value {
	raw, ok := argJSON(args)
	if !ok {
		return jsErr(fmt.Errorf("expected JSON input object"))
	}
	var in struct {
		Spec struct {
			ResistanceOhms float64 `json:"resistanceOhms"`
			PowerWatts     float64 `json:"powerWatts"`
			TolerancePct   float64 `json:"tolerancePct"`
			TempCoeffPPM   int     `json:"tempCoeffPPM"`
			Package        string  `json:"package"`
			Type           string  `json:"type"`
		} `json:"spec"`
		AppliedVoltage float64 `json:"appliedVoltage"`
		AppliedCurrent float64 `json:"appliedCurrent"`
	}
	if err := json.Unmarshal([]byte(raw), &in); err != nil {
		return jsErr(fmt.Errorf("invalid input: %w", err))
	}
	pkg, _ := resistor.ParsePackageType(in.Spec.Package)
	input := resistor.AnalysisInput{
		Spec: resistor.ResistorSpec{
			ResistanceOhms: in.Spec.ResistanceOhms,
			PowerWatts:     in.Spec.PowerWatts,
			TolerancePct:   in.Spec.TolerancePct,
			TempCoeffPPM:   in.Spec.TempCoeffPPM,
			Package:        pkg,
			Type:           resistor.ResistorType(in.Spec.Type),
		},
		AppliedVoltage: in.AppliedVoltage,
		AppliedCurrent: in.AppliedCurrent,
	}
	report, err := resistor.AnalyzeResistor(input)
	if err != nil {
		return jsErr(err)
	}
	return jsOK(report)
}
