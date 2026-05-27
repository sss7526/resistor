package app

import (
	"fmt"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"

	"github.com/sss7526/resistor"
)

// smdInputs is a snapshot of all SMDView inputs for memoizing computeResult.
type smdInputs struct {
	mode       string
	marking    string
	resistance string
	encodeMode resistor.SMDEncodingMode
}

///////////////////////////////////////////////////////////////////////////////
// SMDView
///////////////////////////////////////////////////////////////////////////////

// SMDView provides an interactive form for SMD resistor decode and encode.
//
// Decode mode: SMD marking → resistance value
// Encode mode: resistance value → SMD marking
//
// The Mode field drives a structural form rebuild so the correct input
// fields are always shown. ESC returns to the main menu.
type SMDView struct {
	BaseView

	form *huh.Form

	// Mode tracking for structural rebuild
	mode     string
	prevMode string

	// Decode inputs
	marking string

	// Encode inputs
	resistance string
	encodeMode resistor.SMDEncodingMode

	// Memoization: skip recompute when inputs are unchanged.
	snapshot smdInputs

	// Results — only one is set at a time depending on mode
	decodeResult *resistor.ResistorSpec
	encodeResult string
	err          error
}

///////////////////////////////////////////////////////////////////////////////
// Constructor
///////////////////////////////////////////////////////////////////////////////

func NewSMDView() *SMDView {
	v := &SMDView{
		mode:       "Decode",
		encodeMode: resistor.SMDAuto,
	}
	v.buildForm()
	return v
}

///////////////////////////////////////////////////////////////////////////////
// Form Builder
///////////////////////////////////////////////////////////////////////////////

// buildForm constructs a fresh huh form whose input fields match the current
// mode. Called on construction, on mode change, and on StateCompleted.
func (v *SMDView) buildForm() {
	v.prevMode = v.mode

	modeField := huh.NewSelect[string]().
		Title("Mode").
		Options(
			huh.NewOption("Decode  (marking → Ω)", "Decode"),
			huh.NewOption("Encode  (Ω → marking)", "Encode"),
		).
		Value(&v.mode)

	var fields []huh.Field
	fields = append(fields, modeField)

	if v.mode == "Decode" {
		fields = append(fields,
			huh.NewInput().
				Title("SMD Marking").
				Description("3-digit, 4-digit, R-notation, or EIA-96").
				Value(&v.marking),
		)
	} else {
		fields = append(fields,
			huh.NewInput().
				Title("Resistance (Ω)").
				Value(&v.resistance),
			huh.NewSelect[resistor.SMDEncodingMode]().
				Title("Encoding").
				Options(
					huh.NewOption("Auto (3/4-digit)", resistor.SMDAuto),
					huh.NewOption("Standard (3/4-digit only)", resistor.SMDStandard),
					huh.NewOption("EIA-96", resistor.SMDEIA96),
				).
				Value(&v.encodeMode),
		)
	}

	v.form = huh.NewForm(huh.NewGroup(fields...))

	if v.width > 0 {
		v.form = v.form.WithWidth(v.width/2 - 2)
	}

	// Clear stale results whenever the form is rebuilt.
	v.decodeResult = nil
	v.encodeResult = ""
	v.err = nil
}

///////////////////////////////////////////////////////////////////////////////
// Resize
///////////////////////////////////////////////////////////////////////////////

func (v *SMDView) Resize(width, height int) {
	v.BaseView.Resize(width, height)
	v.form = v.form.WithWidth(width/2 - 2)
}

///////////////////////////////////////////////////////////////////////////////
// Init
///////////////////////////////////////////////////////////////////////////////

func (v *SMDView) Init() tea.Cmd {
	return v.form.Init()
}

///////////////////////////////////////////////////////////////////////////////
// Update
///////////////////////////////////////////////////////////////////////////////

func (v *SMDView) Update(msg tea.Msg) (View, tea.Cmd) {
	updated, cmd := v.form.Update(msg)
	v.form = updated.(*huh.Form)

	// Rebuild on form completion so the panel never goes blank.
	if v.form.State == huh.StateCompleted {
		v.buildForm()
		return v, v.form.Init()
	}

	// ESC navigates back to the menu after huh has had its turn.
	if key, ok := msg.(tea.KeyMsg); ok && key.String() == "esc" {
		return NewMenu(), nil
	}

	// Mode changed: rebuild input fields to match the new selection.
	// huh writes through the mode pointer on every Up/Down keypress in the
	// Select — defer the structural rebuild until the user confirms with
	// Enter or Tab so mid-navigation arrow presses don't clear other fields.
	if v.mode != v.prevMode {
		if key, ok := msg.(tea.KeyMsg); ok {
			switch key.String() {
			case "enter", "tab", "shift+tab":
				v.buildForm()
				return v, v.form.Init()
			}
		}
	}

	v.computeResult()

	return v, cmd
}

///////////////////////////////////////////////////////////////////////////////
// Reactive Computation
///////////////////////////////////////////////////////////////////////////////

func (v *SMDView) computeResult() {
	snap := smdInputs{v.mode, v.marking, v.resistance, v.encodeMode}
	if snap == v.snapshot {
		return
	}
	v.snapshot = snap

	v.err = nil
	v.decodeResult = nil
	v.encodeResult = ""

	if v.mode == "Decode" {
		if v.marking == "" {
			return
		}
		spec, err := resistor.DecodeSMD(v.marking)
		if err != nil {
			v.err = err
			return
		}
		v.decodeResult = &spec

	} else {
		if v.resistance == "" {
			return
		}
		r, err := strconv.ParseFloat(v.resistance, 64)
		if err != nil || r <= 0 {
			v.err = fmt.Errorf("resistance must be a positive number")
			return
		}
		marking, err := resistor.EncodeSMD(r, v.encodeMode)
		if err != nil {
			v.err = err
			return
		}
		v.encodeResult = marking
	}
}

///////////////////////////////////////////////////////////////////////////////
// View Rendering
///////////////////////////////////////////////////////////////////////////////

func (v *SMDView) View() string {
	if v.width <= 0 {
		return ""
	}
	return splitLayout(v.width, v.form.View(), v.renderResult())
}

///////////////////////////////////////////////////////////////////////////////
// Helpers
///////////////////////////////////////////////////////////////////////////////

func (v *SMDView) renderResult() string {
	if v.err != nil {
		return errStyle.Render(v.err.Error())
	}

	if v.mode == "Decode" {
		if v.decodeResult == nil {
			return "Enter an SMD marking to decode."
		}
		return fmt.Sprintf("Resistance: %.6g Ω\n", v.decodeResult.ResistanceOhms)
	}

	if v.encodeResult == "" {
		return "Enter a resistance to encode."
	}
	return fmt.Sprintf("SMD Marking: %s\n", v.encodeResult)
}
