package main

import (
	"fmt"
	"os"

	"github.com/N95Ryan/leaf/internal/app"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	m := app.NewModel()

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
