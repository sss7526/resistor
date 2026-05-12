package app

import (
	tea "github.com/charmbracelet/bubbletea"
)

// View represents a screen in the TUI.
type View interface {
	Init() tea.Cmd
	Update(tea.Msg) (View, tea.Cmd)
	View() string
}