package main

import (
	"context"
	"dailies-go/db"
	"dailies-go/views"
	"fmt"
	"os"

	_ "modernc.org/sqlite"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {

	// Read db
	ctx := context.Background()
	queries, err := db.InitJournalManager("dailies.sqlite")
	if err != nil {
		fmt.Println("Error opening database connection: ", err)
		return
	}

	// Create a new program instance with the HomeView as the initial view
	// WithAltScreen enables full screen mode
	p := tea.NewProgram(views.NewHomeView(queries, ctx), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
