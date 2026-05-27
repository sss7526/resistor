package resistor

import "testing"

///////////////////////////////////////////////////////////////////////////////
// NearestStandard
///////////////////////////////////////////////////////////////////////////////

func BenchmarkNearestStandard_E24(b *testing.B) {
	for b.Loop() {
		_, _ = NearestStandard(487, E24, RoundNearest)
	}
}

func BenchmarkNearestStandard_E96(b *testing.B) {
	for b.Loop() {
		_, _ = NearestStandard(487, E96, RoundNearest)
	}
}

func BenchmarkNearestStandard_E12_RoundUp(b *testing.B) {
	for b.Loop() {
		_, _ = NearestStandard(500, E12, RoundUp)
	}
}

///////////////////////////////////////////////////////////////////////////////
// DecodeBands / EncodeBands
///////////////////////////////////////////////////////////////////////////////

var benchBands4 = []Color{Green, Brown, Brown, Gold}                  // 510Ω ±5%
var benchBands5 = []Color{Yellow, Violet, Black, Red, Brown}          // 4700Ω ±1%
var benchBands6 = []Color{Brown, Black, Black, Orange, Brown, Orange} // 100kΩ ±1% 15ppm

func BenchmarkDecodeBands_4band(b *testing.B) {
	for b.Loop() {
		_, _ = DecodeBands(benchBands4)
	}
}

func BenchmarkDecodeBands_5band(b *testing.B) {
	for b.Loop() {
		_, _ = DecodeBands(benchBands5)
	}
}

func BenchmarkDecodeBands_6band(b *testing.B) {
	for b.Loop() {
		_, _ = DecodeBands(benchBands6)
	}
}

func BenchmarkEncodeBands_4band(b *testing.B) {
	spec := ResistorSpec{ResistanceOhms: 510, TolerancePct: 5}
	for b.Loop() {
		_, _ = EncodeBands(spec)
	}
}

func BenchmarkEncodeBands_5band(b *testing.B) {
	spec := ResistorSpec{ResistanceOhms: 4700, TolerancePct: 1}
	for b.Loop() {
		_, _ = EncodeBands(spec)
	}
}

///////////////////////////////////////////////////////////////////////////////
// SelectStandardResistor
///////////////////////////////////////////////////////////////////////////////

func BenchmarkSelectStandardResistor_E24(b *testing.B) {
	req := SelectionRequest{Resistance: 487, Series: E24, Rounding: RoundNearest}
	for b.Loop() {
		_, _ = SelectStandardResistor(req)
	}
}

func BenchmarkSelectStandardResistor_E96(b *testing.B) {
	req := SelectionRequest{Resistance: 487, Series: E96, Rounding: RoundNearest}
	for b.Loop() {
		_, _ = SelectStandardResistor(req)
	}
}

///////////////////////////////////////////////////////////////////////////////
// InferResistor
///////////////////////////////////////////////////////////////////////////////

func BenchmarkInferResistor_Deterministic(b *testing.B) {
	obs := ObservedResistor{
		Bands: []Color{Yellow, Violet, Brown, Gold},
	}
	for b.Loop() {
		_, _ = InferResistor(obs)
	}
}

func BenchmarkInferResistor_HeuristicOnly(b *testing.B) {
	obs := ObservedResistor{
		BodyColor: Blue,
		LengthMM:  6.3,
	}
	for b.Loop() {
		_, _ = InferResistor(obs)
	}
}

func BenchmarkInferResistor_Combined(b *testing.B) {
	obs := ObservedResistor{
		Bands:     []Color{Yellow, Violet, Brown, Gold},
		BodyColor: Blue,
		LengthMM:  6.3,
		Package:   ThroughHole,
	}
	for b.Loop() {
		_, _ = InferResistor(obs)
	}
}

///////////////////////////////////////////////////////////////////////////////
// DecodeSMD / EncodeSMD
///////////////////////////////////////////////////////////////////////////////

func BenchmarkDecodeSMD_3digit(b *testing.B) {
	for b.Loop() {
		_, _ = DecodeSMD("472")
	}
}

func BenchmarkDecodeSMD_EIA96(b *testing.B) {
	for b.Loop() {
		_, _ = DecodeSMD("01C")
	}
}

func BenchmarkEncodeSMD_Standard(b *testing.B) {
	for b.Loop() {
		_, _ = EncodeSMD(4700, SMDAuto)
	}
}

func BenchmarkEncodeSMD_EIA96(b *testing.B) {
	for b.Loop() {
		_, _ = EncodeSMD(4990, SMDEIA96)
	}
}
