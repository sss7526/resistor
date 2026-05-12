package app

import "github.com/charmbracelet/lipgloss"

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
