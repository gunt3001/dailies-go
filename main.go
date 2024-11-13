package main

import (
	"dailies-go/views"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Create a new program instance with the HomeView as the initial view
	// WithAltScreen enables full screen mode
	p := tea.NewProgram(views.NewHomeView(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
