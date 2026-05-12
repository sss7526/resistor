package app

import (
	tea "github.com/charmbracelet/bubbletea"
)

type viewState int

const (
	viewMenu viewState = iota
	viewSelect
	viewInfer
	viewAnalyze
	viewSMD
	viewQuit
)

// AppModel is the application shell and view router.
//
// It does not implement any business logic.
// Its responsibilities are limited to:
//
//   - Holding the active View
//   - Delegating Update and View calls
//   - Handling global quit keybindings (q, ctrl+c)
//   - Tracking terminal dimensions
//   - Propagating resize events
//   - Applying global layout styling
//
// AppModel intentionally avoids:
//
//   - Routing switches on enum states
//   - View-specific logic
//   - Direct manipulation of child state
//
// View transitions occur when the active view
// returns a new View from its Update method.
type AppModel struct {
	current View
	styles  styles

	width  int
	height int
}

func New() AppModel {
	menu := NewMenu()

	return AppModel{
		current: menu,
		styles:  newStyles(),
	}
}

func (m AppModel) Init() tea.Cmd {
	return m.current.Init()
}

// Update delegates messages to the current view.
//
// Global keys (q, ctrl+c) are handled here.
// Resize events are stored and propagated to views.
//
// When a view transition occurs, the new view is
// immediately resized using the last known dimensions.
func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		if r, ok := m.current.(Resizable); ok {
			r.Resize(m.width, m.height)
		}
	}

	prev := m.current
	next, cmd := m.current.Update(msg)

	// // Always replace current view
	// m.current = next

	// // Always resize after view replacement
	// if m.width > 0 && m.height > 0 {
	// 	if r, ok := m.current.(Resizable); ok {
	// 		r.Resize(m.width, m.height)
	// 	}
	// }

	// If view instance changed, run its Init()
    if next != prev {
        m.current = next

        initCmd := m.current.Init()

        // Immediately resize new view
        if r, ok := m.current.(Resizable); ok {
            r.Resize(m.width, m.height)
        }

        return m, tea.Batch(cmd, initCmd)
    }

    m.current = next

	return m, cmd
}

// View renders the active view wrapped in global layout.
//
// The layout is responsible for header/footer styling.
// Individual views render body content only.
func (m AppModel) View() string {

	body := m.current.View()
	return m.styles.layout(body, m.width)
}
