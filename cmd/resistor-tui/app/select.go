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

	return v
}

///////////////////////////////////////////////////////////////////////////////
// Resize
///////////////////////////////////////////////////////////////////////////////

func (v *SelectView) Resize(width, height int) {
	v.BaseView.Resize(width, height)

	leftWidth := width / 2

	v.form = v.form.WithWidth(leftWidth - 2)
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
	switch msg := msg.(type) {

	case tea.KeyMsg:

		if msg.String() == "esc" {
			return NewMenu(), nil
		}
	}

	var cmd tea.Cmd

	updated, cmd := v.form.Update(msg)
	v.form = updated.(*huh.Form)

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

	formView := v.form.View()

	resultView := v.renderResult()

	return splitLayout(v.width, formView, resultView)
}

///////////////////////////////////////////////////////////////////////////////
// Helpers
///////////////////////////////////////////////////////////////////////////////

func (v *SelectView) renderResult() string {

    if v.err != nil {
        return lipgloss.NewStyle().
            Foreground(lipgloss.Color("#FF5555")).
            Render(v.err.Error())
    }

    if v.result.SelectedResistance == 0 {
        return "Enter values to compute result."
    }

    builder := strings.Builder{}

    builder.WriteString(fmt.Sprintf(
        "Selected: %.6gΩ\n\n",
        v.result.SelectedResistance,
    ))

    builder.WriteString("Bands:\n")
    builder.WriteString(formatBands(v.result.Bands))
    builder.WriteString("\n")

    if len(v.result.Assumptions) > 0 {
        builder.WriteString("Assumptions:\n")
        for _, a := range v.result.Assumptions {
            builder.WriteString(fmt.Sprintf("  - %s\n", a))
        }
    }

    return builder.String()
}

func enumOptions[T interface {
	fmt.Stringer
	comparable
}](values []T) []huh.Option[T] {

	opts := make([]huh.Option[T], len(values))
	for i, v := range values {
		opts[i] = huh.NewOption(v.String(), v)
	}

	return opts
}

func formatBands(bands []resistor.Color) string {
	out := ""
	for _, b := range bands {
		out += fmt.Sprintf("  %s\n", b)
	}
	return out
}
