package app

import (
	"fmt"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"

	"github.com/sss7526/resistor"
)

///////////////////////////////////////////////////////////////////////////////
// AnalyzeView
///////////////////////////////////////////////////////////////////////////////

// AnalyzeView provides an interactive form for analyzing a resistor under
// specified electrical conditions (voltage, current, power rating, tolerance).
//
// The form collects:
//   - Resistance (required)
//   - Applied Voltage (optional)
//   - Applied Current (optional)
//   - Rated Power (optional; enables derating analysis)
//   - Tolerance % (optional; enables worst-case bounds)
//
// Results are computed reactively after each form update.
// ESC returns to the main menu.
type AnalyzeView struct {
	BaseView

	form *huh.Form

	// Bound variables
	resistance string
	voltage    string
	current    string
	power      string
	tolerance  string

	// Computed result
	result    resistor.AnalysisReport
	hasResult bool
	err       error
}

///////////////////////////////////////////////////////////////////////////////
// Constructor
///////////////////////////////////////////////////////////////////////////////

func NewAnalyzeView() *AnalyzeView {
	v := &AnalyzeView{}

	v.form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Resistance (Ω)").
				Value(&v.resistance),

			huh.NewInput().
				Title("Applied Voltage (V)").
				Value(&v.voltage),

			huh.NewInput().
				Title("Applied Current (A)").
				Value(&v.current),

			huh.NewInput().
				Title("Rated Power (W)").
				Value(&v.power),

			huh.NewInput().
				Title("Tolerance (%)").
				Value(&v.tolerance),
		),
	)

	return v
}

///////////////////////////////////////////////////////////////////////////////
// Resize
///////////////////////////////////////////////////////////////////////////////

func (v *AnalyzeView) Resize(width, height int) {
	v.BaseView.Resize(width, height)
	leftWidth := width / 2
	v.form = v.form.WithWidth(leftWidth - 2)
}

///////////////////////////////////////////////////////////////////////////////
// Init
///////////////////////////////////////////////////////////////////////////////

func (v *AnalyzeView) Init() tea.Cmd {
	return v.form.Init()
}

///////////////////////////////////////////////////////////////////////////////
// Update
///////////////////////////////////////////////////////////////////////////////

func (v *AnalyzeView) Update(msg tea.Msg) (View, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "esc" {
			return NewMenu(), nil
		}
	}

	updated, cmd := v.form.Update(msg)
	v.form = updated.(*huh.Form)

	v.computeResult()

	return v, cmd
}

///////////////////////////////////////////////////////////////////////////////
// Reactive Computation
///////////////////////////////////////////////////////////////////////////////

func (v *AnalyzeView) computeResult() {
	if v.resistance == "" {
		v.err = nil
		v.hasResult = false
		return
	}

	r, err := strconv.ParseFloat(v.resistance, 64)
	if err != nil || r <= 0 {
		v.err = fmt.Errorf("invalid resistance")
		v.hasResult = false
		return
	}

	spec := resistor.ResistorSpec{ResistanceOhms: r}

	if v.power != "" {
		pw, err := strconv.ParseFloat(v.power, 64)
		if err == nil && pw > 0 {
			spec.PowerWatts = pw
		}
	}

	if v.tolerance != "" {
		tol, err := strconv.ParseFloat(v.tolerance, 64)
		if err == nil && tol > 0 {
			spec.TolerancePct = tol
		}
	}

	input := resistor.AnalysisInput{Spec: spec}

	if v.voltage != "" {
		vv, err := strconv.ParseFloat(v.voltage, 64)
		if err == nil {
			input.AppliedVoltage = vv
		}
	}

	if v.current != "" {
		ic, err := strconv.ParseFloat(v.current, 64)
		if err == nil {
			input.AppliedCurrent = ic
		}
	}

	result, err := resistor.AnalyzeResistor(input)
	if err != nil {
		v.err = err
		v.hasResult = false
		return
	}

	v.result = result
	v.hasResult = true
	v.err = nil
}

///////////////////////////////////////////////////////////////////////////////
// View Rendering
///////////////////////////////////////////////////////////////////////////////

func (v *AnalyzeView) View() string {
	return splitLayout(v.width, v.form.View(), v.renderResult())
}

///////////////////////////////////////////////////////////////////////////////
// Helpers
///////////////////////////////////////////////////////////////////////////////

func (v *AnalyzeView) renderResult() string {
	if v.err != nil {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5555")).
			Render(v.err.Error())
	}

	if !v.hasResult {
		return "Enter resistance to compute analysis."
	}

	rep := v.result
	var b strings.Builder

	fmt.Fprintf(&b, "Power Dissipation: %.4g W\n", rep.PowerDissipation)
	fmt.Fprintf(&b, "Voltage Drop:      %.4g V\n", rep.VoltageDrop)
	fmt.Fprintf(&b, "Current:           %.4g A\n\n", rep.Current)

	if rep.DeratedSafePower != nil {
		fmt.Fprintf(&b, "Derated Safe Power (50%%): %.4g W\n\n", *rep.DeratedSafePower)
	}

	if rep.WorstCaseResistanceMin != nil && rep.WorstCaseResistanceMax != nil {
		fmt.Fprintf(&b, "Worst-Case Resistance:\n")
		fmt.Fprintf(&b, "  Min: %.6g Ω\n", *rep.WorstCaseResistanceMin)
		fmt.Fprintf(&b, "  Max: %.6g Ω\n\n", *rep.WorstCaseResistanceMax)
	}

	b.WriteString(renderWarnings(rep.Warnings))

	return b.String()
}
