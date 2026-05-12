package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"log"

	"github.com/sss7526/resistor/cmd/resistor-tui/app"
)

func main() {
	p := tea.NewProgram(app.New(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
