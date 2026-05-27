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
type SelectView struct {
	BaseView

	form *huh.Form

	// Bound variables
	resistance string
	tolerance  string
	series     resistor.ESeries
	rounding   resistor.RoundingMode

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
	updated, cmd := v.form.Update(msg)
	v.form = updated.(*huh.Form)

	if v.form.State == huh.StateCompleted {
		v.buildForm()
		return v, v.form.Init()
	}

	if key, ok := msg.(tea.KeyMsg); ok && key.String() == "esc" {
		return NewMenu(), nil
	}

	v.computeResult()

	return v, cmd
}

///////////////////////////////////////////////////////////////////////////////
// Reactive Computation
///////////////////////////////////////////////////////////////////////////////

func (v *SelectView) computeResult() {

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
		tol, _ = strconv.ParseFloat(v.tolerance, 64)
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

	builder := strings.Builder{}

	fmt.Fprintf(&builder, "Selected: %.6gΩ\n\n", v.result.SelectedResistance)

	builder.WriteString("Bands:\n")
	builder.WriteString(formatBands(v.result.Bands))
	builder.WriteString("\n")

	if len(v.result.Assumptions) > 0 {
		builder.WriteString("Assumptions:\n")
		for _, a := range v.result.Assumptions {
			fmt.Fprintf(&builder, "  - %s\n", a)
		}
	}

	return builder.String()
}

func formatBands(bands []resistor.Color) string {
	var b strings.Builder
	for _, c := range bands {
		fmt.Fprintf(&b, "  %s\n", c)
	}
	return b.String()
}
