package app

import (
	"github.com/charmbracelet/lipgloss"
)

type styles struct {
	header lipgloss.Style
	panel  lipgloss.Style
	footer lipgloss.Style
}

// styles defines the global visual theme for the TUI.
//
// All layout styling is centralized here to ensure:
//
//   - Consistent look and feel
//   - No inline styling scattered across views
//   - Easy future theming adjustments
//
// Individual views render body content only.
// AppModel wraps body content using layout().
func newStyles() styles {
	return styles{

		header: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7D56F4")).
			Padding(0, 1),

		panel: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(1, 2),

		footer: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			Padding(0, 1),
	}
}

// layout composes header, body panel, and footer
// into a vertically stacked layout.
//
// width is provided by AppModel and reflects
// the current terminal width.
func (s styles) layout(body string, width int) string {
	header := s.header.Width(width).Render("Resistor Engineering Toolkit")
	panel := s.panel.Width(width - 4).Render(body)
	footer := s.footer.Width(width).Render("q: quit • esc: back")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		panel,
		footer,
	)
}
