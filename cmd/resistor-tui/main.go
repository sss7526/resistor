package main

import (
	"log"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/sss7526/resistor/cmd/resistor-tui/app"
)

func main() {
	p := tea.NewProgram(app.New(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}