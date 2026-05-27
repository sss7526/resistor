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

	// Input snapshot for memoization: skip recompute when nothing changed.
	snapshot [5]string

	// Computed result; nil means no result has been computed yet.
	result *resistor.AnalysisReport
	err    error
}

///////////////////////////////////////////////////////////////////////////////
// Constructor
///////////////////////////////////////////////////////////////////////////////

func NewAnalyzeView() *AnalyzeView {
	v := &AnalyzeView{}
	v.buildForm()
	return v
}

///////////////////////////////////////////////////////////////////////////////
// Form Builder
///////////////////////////////////////////////////////////////////////////////

// buildForm constructs a fresh huh form bound to the existing field pointers.
// Called on construction and again when the form reaches StateCompleted so the
// user is never left with a blank panel.
func (v *AnalyzeView) buildForm() {
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

	if v.width > 0 {
		v.form = v.form.WithWidth(v.width/2 - 2)
	}
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
	// Pass message to form first so huh can handle any field-level keys
	// (e.g. a future Select field using ESC to close its dropdown).
	updated, cmd := v.form.Update(msg)
	v.form = updated.(*huh.Form)

	// When huh completes the form (user pressed Enter past the last field),
	// rebuild so the panel doesn't go blank and the user can keep editing.
	if v.form.State == huh.StateCompleted {
		v.buildForm()
		return v, v.form.Init()
	}

	// ESC navigates back to the menu after huh has had its turn.
	if key, ok := msg.(tea.KeyMsg); ok && key.String() == "esc" {
		return NewMenu(), nil
	}

	v.computeResult()

	return v, cmd
}

///////////////////////////////////////////////////////////////////////////////
// Reactive Computation
///////////////////////////////////////////////////////////////////////////////

func (v *AnalyzeView) computeResult() {
	// Memoize: skip if no input field has changed since last run.
	snap := [5]string{v.resistance, v.voltage, v.current, v.power, v.tolerance}
	if snap == v.snapshot {
		return
	}
	v.snapshot = snap

	v.result = nil
	v.err = nil

	if v.resistance == "" {
		return
	}

	r, err := strconv.ParseFloat(v.resistance, 64)
	if err != nil || r <= 0 {
		v.err = fmt.Errorf("resistance must be a positive number")
		return
	}

	spec := resistor.ResistorSpec{ResistanceOhms: r}

	// Optional fields: if non-empty they must parse to a positive number.
	if v.power != "" {
		pw, err := strconv.ParseFloat(v.power, 64)
		if err != nil || pw <= 0 {
			v.err = fmt.Errorf("rated power must be a positive number")
			return
		}
		spec.PowerWatts = pw
	}

	if v.tolerance != "" {
		tol, err := strconv.ParseFloat(v.tolerance, 64)
		if err != nil || tol <= 0 {
			v.err = fmt.Errorf("tolerance must be a positive number")
			return
		}
		spec.TolerancePct = tol
	}

	input := resistor.AnalysisInput{Spec: spec}

	if v.voltage != "" {
		vv, err := strconv.ParseFloat(v.voltage, 64)
		if err != nil || vv <= 0 {
			v.err = fmt.Errorf("voltage must be a positive number")
			return
		}
		input.AppliedVoltage = vv
	}

	if v.current != "" {
		ic, err := strconv.ParseFloat(v.current, 64)
		if err != nil || ic <= 0 {
			v.err = fmt.Errorf("current must be a positive number")
			return
		}
		input.AppliedCurrent = ic
	}

	report, err := resistor.AnalyzeResistor(input)
	if err != nil {
		v.err = err
		return
	}

	v.result = &report
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

	if v.result == nil {
		return "Enter resistance to compute analysis."
	}

	rep := *v.result
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
