package app

import (
	tea "github.com/charmbracelet/bubbletea"
)

// View represents a screen in the TUI.
// Each view is autonomouse and may return a new View during Update.
type View interface {
	Init() tea.Cmd
	Update(tea.Msg) (View, tea.Cmd)
	View() string
}

// Resizable is implmemented by views that need to react to terminal resizing.
type Resizable interface {
	Resize(width, height int)
}
