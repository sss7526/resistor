package app

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"

	"github.com/sss7526/resistor"
)

///////////////////////////////////////////////////////////////////////////////
// InferView
///////////////////////////////////////////////////////////////////////////////

/*
InferView provides structured input for resistor inference.

Users may select:

  - Input mode: Bands or SMD
  - Band count (4/5/6)
  - Individual band colors
  - Body color
  - Length
  - Package type

The view reactively calls resistor.InferResistor and renders:

  - Electrical properties
  - Inferred type
  - Power rating
  - Voltage rating
  - Confidence
  - Assumptions

ESC returns to the main menu.
*/
type InferView struct {
	BaseView

	form *huh.Form

	coreGroup     *huh.Group
	physicalGroup *huh.Group

	viewport viewport.Model

	mode string

	bandCount int
	bands     []resistor.Color
	smd       string

	bodyColor resistor.Color
	length    string
	pkg       resistor.PackageType

	result resistor.InferenceResult
	err    error
}

///////////////////////////////////////////////////////////////////////////////
// Constructor
///////////////////////////////////////////////////////////////////////////////

func NewInferView() *InferView {

	v := &InferView{
		mode:      "Bands",
		bandCount: 4,
		bands:     make([]resistor.Color, 6),
	}
	v.viewport = viewport.New(0, 0)
	v.viewport.SetContent("")
	v.buildForm()

	return v
}

///////////////////////////////////////////////////////////////////////////////
// Resize
///////////////////////////////////////////////////////////////////////////////

func (v *InferView) Resize(width, height int) {
	v.BaseView.Resize(width, height)

	totalWidth := width
	formWidth := totalWidth / 2

	// Viewport only occupies left half
	v.viewport.Width = formWidth
	v.viewport.Height = height - 4
}

///////////////////////////////////////////////////////////////////////////////
// Init
///////////////////////////////////////////////////////////////////////////////

func (v *InferView) Init() tea.Cmd {
	return v.form.Init()
}

///////////////////////////////////////////////////////////////////////////////
// Form Builder
///////////////////////////////////////////////////////////////////////////////

func (v *InferView) buildForm() {

	// ----- Core Inputs Group -----
	coreFields := []huh.Field{
		huh.NewSelect[string]().
			Title("Input Mode").
			Options(
				huh.NewOption("Bands", "Bands"),
				huh.NewOption("SMD", "SMD"),
			).
			Value(&v.mode),
	}

	if v.mode == "Bands" {

		coreFields = append(coreFields,
			huh.NewSelect[int]().
				Title("Band Count").
				Options(
					huh.NewOption("4", 4),
					huh.NewOption("5", 5),
					huh.NewOption("6", 6),
				).
				Value(&v.bandCount),
		)

		for i := 0; i < v.bandCount; i++ {
			coreFields = append(coreFields,
				huh.NewSelect[resistor.Color]().
					Title(fmt.Sprintf("Band %d", i+1)).
					Options(enumOptions(resistor.DigitColors())...).
					Value(&v.bands[i]),
			)
		}

	} else {
		coreFields = append(coreFields,
			huh.NewInput().
				Title("SMD Marking").
				Value(&v.smd),
		)
	}

	v.coreGroup = huh.NewGroup(coreFields...)

	// ----- Physical Properties Group -----
	physicalFields := []huh.Field{
		huh.NewSelect[resistor.Color]().
			Title("Body Color").
			Options(enumOptions(resistor.DigitColors())...).
			Value(&v.bodyColor),

		huh.NewInput().
			Title("Length (mm)").
			Value(&v.length),

		huh.NewSelect[resistor.PackageType]().
			Title("Package").
			Options(enumOptions(resistor.AllPackageTypes())...).
			Value(&v.pkg),
	}

	v.physicalGroup = huh.NewGroup(physicalFields...)

	v.form = huh.NewForm(
		v.coreGroup,
		v.physicalGroup,
	)
}

///////////////////////////////////////////////////////////////////////////////
// Update
///////////////////////////////////////////////////////////////////////////////

func (v *InferView) Update(msg tea.Msg) (View, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "esc" {
			return NewMenu(), nil
		}
	}

	updated, cmd := v.form.Update(msg)
	v.form = updated.(*huh.Form)

	v.viewport, _ = v.viewport.Update(msg)

	v.computeResult()

	return v, cmd
}

///////////////////////////////////////////////////////////////////////////////
// Reactive Computation
///////////////////////////////////////////////////////////////////////////////

func (v *InferView) computeResult() {

	obs := resistor.ObservedResistor{}

	if v.mode == "Bands" {
		obs.Bands = v.bands[:v.bandCount]
	} else {
		obs.Marking = v.smd
	}

	if v.length != "" {
		lengthVal, err := strconv.ParseFloat(v.length, 64)
		if err == nil {
			obs.LengthMM = lengthVal
		}
	}

	obs.BodyColor = v.bodyColor
	obs.Package = v.pkg

	result, err := resistor.InferResistor(obs)
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

func (v *InferView) View() string {

	if v.width <= 0 {
		return ""
	}

	totalWidth := v.width
	formWidth := totalWidth / 2
	resultWidth := totalWidth - formWidth - 2

	coreWidth := formWidth / 2
	physicalWidth := formWidth - coreWidth - 2

	coreView := lipgloss.NewStyle().
		Width(coreWidth).
		Render(v.coreGroup.View())

	physicalView := lipgloss.NewStyle().
		Width(physicalWidth).
		Render(v.physicalGroup.View())

	formCombined := lipgloss.JoinHorizontal(
		lipgloss.Top,
		coreView,
		"  ",
		physicalView,
	)

	// Set viewport content
	v.viewport.SetContent(formCombined)

	resultPanel := lipgloss.NewStyle().
		Width(resultWidth).
		Render(v.renderResult())

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		v.viewport.View(),
		"  ",
		resultPanel,
	)
}

func (v *InferView) renderResult() string {

	if v.err != nil {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5555")).
			Render(v.err.Error())
	}

	if v.result.Spec.ResistanceOhms == 0 {
		return "Enter values to compute inference."
	}

	out := fmt.Sprintf(
		"Resistance: %.6gΩ\nTolerance: %.2f%%\n",
		v.result.Spec.ResistanceOhms,
		v.result.Spec.TolerancePct,
	)

	out += fmt.Sprintf(
		"Power: %.3gW\nType: %s\nVoltage: %.3gV\n\n",
		v.result.Spec.PowerWatts,
		v.result.Spec.Type,
		v.result.VoltageRating,
	)

	out += renderConfidence(v.result.Meta.Confidence)
	out += "\n"
	out += renderAssumptions(v.result.Meta.Assumptions)

	return out
}
