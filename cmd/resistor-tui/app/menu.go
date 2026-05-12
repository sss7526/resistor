package app

import (
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/bubbles/list"
)

type menuItem struct {
    title string
    desc  string
}

func (i menuItem) Title() string       { return i.title }
func (i menuItem) Description() string { return i.desc }
func (i menuItem) FilterValue() string { return i.title }

type MenuView struct {
    list list.Model
}

func NewMenu() MenuView {
    items := []list.Item{
        menuItem{"Select Resistor", "Snap to standard values"},
        menuItem{"Infer Resistor", "Analyze from physical clues"},
        menuItem{"Analyze Resistor", "Electrical power & safety"},
        menuItem{"SMD Tools", "Encode / decode SMD markings"},
        menuItem{"Quit", "Exit the application"},
    }

    l := list.New(items, list.NewDefaultDelegate(), 0, 0)
    l.Title = "Main Menu"
    l.SetShowHelp(true)

    return MenuView{list: l}
}

func (m MenuView) Init() tea.Cmd {
    return nil
}

func (m MenuView) Update(msg tea.Msg) (View, tea.Cmd) {

    switch msg := msg.(type) {

    case tea.WindowSizeMsg:
        m.list.SetSize(msg.Width-4, msg.Height-6)

    case tea.KeyMsg:
    	switch msg.String() {
		case "esc":
			// Only block ESC if NOT filtering
			if m.list.FilterState() == list.Unfiltered {
				return m, nil
			}
			// Otherwise, let list handle ESC
		case "enter":
            switch m.list.Index() {
            case 0:
                return NewPlaceholderView("Select View (not yet implemented)"), nil
            case 1:
                return NewPlaceholderView("Infer View (not yet implemented)"), nil
            case 2:
                return NewPlaceholderView("Analyze View (not yet implemented)"), nil
            case 3:
                return NewPlaceholderView("SMD View (not yet implemented)"), nil
            case 4:
                return m, tea.Quit
            }
        }
    }

    var cmd tea.Cmd
    m.list, cmd = m.list.Update(msg)
    return m, cmd
}

func (m MenuView) View() string {
    return m.list.View()
}