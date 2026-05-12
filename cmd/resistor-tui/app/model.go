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

type AppModel struct {
	current View
	styles styles

	width int
	height int
}

func New() AppModel {
	menu := NewMenu()

	return AppModel{
		current: menu,
		styles: newStyles(),
	}
}

func (m AppModel) Init() tea.Cmd {
	return m.current.Init()
}

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
    }

    next, cmd := m.current.Update(msg)
    m.current = next

    if m.width > 0 && m.height > 0 {
        resizeCmd := func() tea.Msg {
            return tea.WindowSizeMsg{
                Width:  m.width,
                Height: m.height,
            }
        }
        return m, tea.Batch(cmd, resizeCmd)
    }

    return m, cmd
}

func (m AppModel) View() string {

	body := m.current.View()
	return m.styles.layout(body, m.width)
}