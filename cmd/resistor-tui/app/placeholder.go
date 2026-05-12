package app

import (
	tea "github.com/charmbracelet/bubbletea"
)

// PlaceholderView is a temporary stand-in for
// unimplemented views.
//
// It exists to allow incremental development
// without breaking the routing architecture.
//
// ESC returns to the main menu.
type PlaceholderView struct {
	message string
}

func NewPlaceholderView(msg string) *PlaceholderView {
	return &PlaceholderView{message: msg}
}

func (p *PlaceholderView) Init() tea.Cmd {
	return nil
}

func (p *PlaceholderView) Update(msg tea.Msg) (View, tea.Cmd) {

	if key, ok := msg.(tea.KeyMsg); ok {
		if key.String() == "esc" {
			return NewMenu(), nil
		}
	}

	return p, nil
}

func (p *PlaceholderView) View() string {
	return p.message + "\n\n(Press ESC to return)"
}
