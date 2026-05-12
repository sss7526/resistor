package app

import (
	"github.com/charmbracelet/lipgloss"
)

type styles struct {
	header lipgloss.Style
	panel lipgloss.Style
	footer lipgloss.Style
}

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

func (s styles) layout(body string, width int) string {
	header := s.header.Width(width).Render("Resistor Engineering Toolkit")
	panel := s.panel.Width(width-4).Render(body)
	footer := s.footer.Width(width).Render("q: quit • esc: back")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		panel,
		footer,
	)
}