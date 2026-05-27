package app

import (
	"fmt"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"strings"

	"github.com/sss7526/resistor"
)

// BaseView provides shared width/height tracking
// for views that depend on terminal dimensions.
//
// Views embed BaseView to avoid repeating
// dimension storage boilerplate.
//
// Resize must be implemented by the view
// if it needs to adjust internal components.
type BaseView struct {
	width  int
	height int
}

// Resize stores the current terminal dimensions.
//
// Views embedding BaseView may extend this method
// to adjust internal Bubble components (e.g., list.SetSize).
func (b *BaseView) Resize(width, height int) {
	b.width = width
	b.height = height
}

// splitLayout renders a left/right split panel.
//
// leftWidth is proportional; remaining width goes to right.
func splitLayout(width int, left string, right string) string {
	if width <= 0 {
		return left + "\n" + right
	}

	leftWidth := width / 2
	rightWidth := width - leftWidth - 2

	leftPanel := lipgloss.NewStyle().
		Width(leftWidth).
		Render(left)

	rightPanel := lipgloss.NewStyle().
		Width(rightWidth).
		Render(right)

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftPanel,
		"  ",
		rightPanel,
	)
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

// errStyle is the shared error text style used by all result panels.
var errStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5555"))

func renderBands(bands []resistor.Color) string { //nolint:unused // used by SMD view when implemented
	var b strings.Builder
	for _, c := range bands {
		fmt.Fprintf(&b, "  %s\n", c)
	}
	return b.String()
}

func renderAssumptions(assumptions []string) string {
	if len(assumptions) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString("Assumptions:\n")
	for _, a := range assumptions {
		fmt.Fprintf(&b, "  - %s\n", a)
	}
	return b.String()
}

func renderWarnings(warnings []resistor.AnalysisWarning) string { //nolint:unused // used by Analyze/SMD views when implemented
	if len(warnings) == 0 {
		return ""
	}

	var b strings.Builder
	b.WriteString("Warnings:\n")

	for _, w := range warnings {
		style := lipgloss.NewStyle()

		switch w.Level {
		case resistor.WarningDanger:
			style = style.Foreground(lipgloss.Color("#FF0000"))
		case resistor.WarningCaution:
			style = style.Foreground(lipgloss.Color("#FFA500"))
		case resistor.WarningInfo:
			style = style.Foreground(lipgloss.Color("#AAAAAA"))
		}

		b.WriteString(style.Render(fmt.Sprintf("  - %s\n", w.Message)))
	}

	return b.String()
}

func renderConfidence(conf float64) string {
	barWidth := 20
	filled := int(conf * float64(barWidth))

	bar := strings.Repeat("█", filled) +
		strings.Repeat("░", barWidth-filled)

	return fmt.Sprintf("Confidence: [%s] %.2f\n", bar, conf)
}
