package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	tui "github.com/corpus/internal/tui/blog"
)

func main() {
	p := tea.NewProgram(tui.InitialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
