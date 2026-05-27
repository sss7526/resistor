package app

import (
	"fmt"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"

	"github.com/sss7526/resistor"
)

///////////////////////////////////////////////////////////////////////////////
// SelectView (Form-Based)
///////////////////////////////////////////////////////////////////////////////

/*
SelectView provides an interactive form for selecting
a standard resistor value.

This view leverages the `huh` form system to manage:

  - Focus navigation
  - Field validation
  - Structured input grouping

The form collects:

  - Resistance (required)
  - Tolerance (optional)
  - E-Series (select)
  - Rounding mode (select)

Results are computed reactively after form updates.

SelectView does not manage routing logic.
ESC returns to the main menu.
*/

// selectInputs holds a snapshot of all SelectView form fields for memoizing computeResult.
// Named struct so the compiler flags any mismatch when fields are added or removed.
type selectInputs struct {
	resistance string
	tolerance  string
	series     resistor.ESeries
	rounding   resistor.RoundingMode
}

type SelectView struct {
	BaseView

	form *huh.Form

	// Bound variables
	resistance string
	tolerance  string
	series     resistor.ESeries
	rounding   resistor.RoundingMode

	// Input snapshot for memoization: skip recompute when nothing changed.
	snapshot selectInputs

	// Computed result
	result resistor.SelectionResult
	err    error
}

///////////////////////////////////////////////////////////////////////////////
// Constructor
///////////////////////////////////////////////////////////////////////////////

func NewSelectView() *SelectView {
	v := &SelectView{}
	v.series = resistor.E24
	v.rounding = resistor.RoundNearest
	v.buildForm()
	return v
}

func (v *SelectView) buildForm() {
	v.form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Resistance (Ω)").
				Value(&v.resistance),

			huh.NewInput().
				Title("Tolerance (%)").
				Value(&v.tolerance),

			huh.NewSelect[resistor.ESeries]().
				Title("E-Series").
				Options(enumOptions(resistor.AllESeries())...).
				Value(&v.series),

			huh.NewSelect[resistor.RoundingMode]().
				Title("Rounding").
				Options(enumOptions(resistor.AllRoundingModes())...).
				Value(&v.rounding),
		),
	)
	if v.width > 0 {
		v.form = v.form.WithWidth(v.width/2 - 2)
	}
}

///////////////////////////////////////////////////////////////////////////////
// Resize
///////////////////////////////////////////////////////////////////////////////

func (v *SelectView) Resize(width, height int) {
	v.BaseView.Resize(width, height)
	v.form = v.form.WithWidth(width/2 - 2)
}

///////////////////////////////////////////////////////////////////////////////
// Init
///////////////////////////////////////////////////////////////////////////////

func (v *SelectView) Init() tea.Cmd {
	return v.form.Init()
}

///////////////////////////////////////////////////////////////////////////////
// Update
///////////////////////////////////////////////////////////////////////////////

func (v *SelectView) Update(msg tea.Msg) (View, tea.Cmd) {
	// ESC checked before form.Update so it always exits the view. If checked
	// after, huh's Select filter mode consumes ESC to clear the filter and
	// the same message also fires the navigation check, ejecting the user.
	if key, ok := msg.(tea.KeyMsg); ok && key.String() == "esc" {
		return NewMenu(), nil
	}

	updated, cmd := v.form.Update(msg)
	v.form = updated.(*huh.Form)

	if v.form.State == huh.StateCompleted {
		v.buildForm()
		return v, v.form.Init()
	}

	v.computeResult()

	return v, cmd
}

///////////////////////////////////////////////////////////////////////////////
// Reactive Computation
///////////////////////////////////////////////////////////////////////////////

func (v *SelectView) computeResult() {
	snap := selectInputs{v.resistance, v.tolerance, v.series, v.rounding}
	if snap == v.snapshot {
		return
	}
	v.snapshot = snap

	if v.resistance == "" {
		v.err = nil
		return
	}

	value, err := strconv.ParseFloat(v.resistance, 64)
	if err != nil || value <= 0 {
		v.err = fmt.Errorf("invalid resistance")
		return
	}

	tol := 0.0
	if v.tolerance != "" {
		var tolErr error
		tol, tolErr = strconv.ParseFloat(v.tolerance, 64)
		if tolErr != nil || tol < 0 || tol >= 100 {
			v.err = fmt.Errorf("tolerance must be between 0 and 100 (exclusive)")
			return
		}
	}

	req := resistor.SelectionRequest{
		Resistance:   value,
		TolerancePct: tol,
		Series:       v.series,
		Rounding:     v.rounding,
	}

	result, err := resistor.SelectStandardResistor(req)
	if err != nil {
		v.err = err
		return
	}
	v.result = result
	v.err = nil
}

///////////////////////////////////////////////////////////////////////////////
// View Rendering
///////////////////////////////////////////////////////////////////////////////

func (v *SelectView) View() string {
	if v.width <= 0 {
		return ""
	}
	return splitLayout(v.width, v.form.View(), v.renderResult())
}

///////////////////////////////////////////////////////////////////////////////
// Helpers
///////////////////////////////////////////////////////////////////////////////

func (v *SelectView) renderResult() string {

	if v.err != nil {
		return errStyle.Render(v.err.Error())
	}

	if v.result.SelectedResistance == 0 {
		return "Enter values to compute result."
	}

	var builder strings.Builder

	fmt.Fprintf(&builder, "Selected: %.6gΩ\n\n", v.result.SelectedResistance)

	builder.WriteString("Bands:\n")
	builder.WriteString(renderBands(v.result.Bands))
	builder.WriteString("\n")

	builder.WriteString(renderAssumptions(v.result.Assumptions))

	return builder.String()
}
