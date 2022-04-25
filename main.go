package main

import (
	"log"

	"github.com/c-grimshaw/gosniff/cmd/gosniff"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(gosniff.NewModel())
	if err := p.Start(); err != nil {
		log.Fatalf("Alas, there's been an error: %v", err)
	}
}
