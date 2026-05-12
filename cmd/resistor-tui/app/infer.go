package app

import (
	"fmt"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"

	"github.com/sss7526/resistor"
)

type InferView struct {
	BaseView

	form          *huh.Form
	coreGroup     *huh.Group
	physicalGroup *huh.Group

	// Mode + structural state
	mode          string
	prevMode      string
	bandCount     int
	prevBandCount int

	// Inputs
	bands     []resistor.Color
	smd       string
	bodyColor resistor.Color
	length    string
	pkg       resistor.PackageType

	// Result
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

	v.buildForm()

	return v
}

///////////////////////////////////////////////////////////////////////////////
// Resize
///////////////////////////////////////////////////////////////////////////////

func (v *InferView) Resize(width, height int) {
	v.BaseView.Resize(width, height)
}

///////////////////////////////////////////////////////////////////////////////
// Init
///////////////////////////////////////////////////////////////////////////////

func (v *InferView) Init() tea.Cmd {
	return v.form.Init()
}

///////////////////////////////////////////////////////////////////////////////
// Form Builder (Role-driven)
///////////////////////////////////////////////////////////////////////////////

func (v *InferView) buildForm() {

	v.prevMode = v.mode
	v.prevBandCount = v.bandCount

	var coreFields []huh.Field

	// Mode selector
	coreFields = append(coreFields,
		huh.NewSelect[string]().
			Title("Input Mode").
			Options(
				huh.NewOption("Bands", "Bands"),
				huh.NewOption("SMD", "SMD"),
			).
			Value(&v.mode),
	)

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

		roles, _ := resistor.BandRolesForCount(v.bandCount)

		for i, role := range roles {

			validColors := resistor.ValidColorsForRole(role)

			coreFields = append(coreFields,
				huh.NewSelect[resistor.Color]().
					Title(fmt.Sprintf("Band %d (%s)", i+1, role.String())).
					Options(enumOptions(validColors)...).
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

	v.physicalGroup = huh.NewGroup(
		huh.NewSelect[resistor.Color]().
			Title("Body Color").
			Options(enumOptions(resistor.BodyColors())...).
			Value(&v.bodyColor),

		huh.NewInput().
			Title("Length (mm)").
			Value(&v.length),

		huh.NewSelect[resistor.PackageType]().
			Title("Package").
			Options(enumOptions(resistor.AllPackageTypes())...).
			Value(&v.pkg),
	)

	v.form = huh.NewForm(v.coreGroup, v.physicalGroup)
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

	// Structural change detection
	if v.mode != v.prevMode || v.bandCount != v.prevBandCount {
		v.buildForm()
		v.Resize(v.width, v.height)
		return v, v.form.Init()
	}

	v.computeResult()

	return v, cmd
}

///////////////////////////////////////////////////////////////////////////////
// Reactive Inference
///////////////////////////////////////////////////////////////////////////////

func (v *InferView) computeResult() {

	obs := resistor.ObservedResistor{}

	if v.mode == "Bands" {
		roles, _ := resistor.BandRolesForCount(v.bandCount)
		obs.Bands = v.bands[:len(roles)]
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
// View Rendering (Split Layout)
///////////////////////////////////////////////////////////////////////////////

func (v *InferView) View() string {

	if v.width <= 0 {
		return ""
	}

	totalWidth := v.width
	formWidth := totalWidth / 2
	resultWidth := totalWidth - formWidth - 2

	// Split form into two columns
	leftWidth := formWidth / 2
	rightWidth := formWidth - leftWidth - 2

	corePanel := lipgloss.NewStyle().
		Width(leftWidth).
		Render(v.coreGroup.View())

	physicalPanel := lipgloss.NewStyle().
		Width(rightWidth).
		Render(v.physicalGroup.View())

	formSplit := lipgloss.JoinHorizontal(
		lipgloss.Top,
		corePanel,
		"  ",
		physicalPanel,
	)

	resultPanel := lipgloss.NewStyle().
		Width(resultWidth).
		Render(v.renderResult())

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		formSplit,
		"  ",
		resultPanel,
	)
}

///////////////////////////////////////////////////////////////////////////////
// Result Rendering
///////////////////////////////////////////////////////////////////////////////

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
