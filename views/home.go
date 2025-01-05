package views

import (
	"context"
	"dailies-go/db"
	"dailies-go/views/models"
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	sidebarProportion = 20 // Width of the sidebar as an integer percent
	helpHeight        = 1  // Height of the help view (in lines)
)

// Struct to handle key bindings
type homeViewKeymap struct {
	quit   key.Binding
	reroll key.Binding
}

// HomeView - Primary view for the application
type HomeView struct {
	width, height, sidebarWidth int
	help                        help.Model
	keymap                      homeViewKeymap
	database                    *db.Queries
	dbContext                   context.Context

	// Styles
	borderStyle   lipgloss.Style
	mainAreaStyle lipgloss.Style
	sidebarStyle  lipgloss.Style

	// Data
	entriesList list.Model
	randomEntry *db.Entry
}

// NewHomeView Constructor
func NewHomeView(database *db.Queries, context context.Context) HomeView {

	// Initialize entries from db
	entries, err := database.GetEntriesByDateRange(context, db.GetEntriesByDateRangeParams{
		DateStart: "2024-01-01",
		DateEnd:   "2024-12-31",
	})
	// TODO: Handle database read errors
	_ = err
	// Convert into list format
	entryListItems := make([]list.Item, len(entries))
	for i := range entries {
		entryListItems[i] = entries[i]
	}
	// We use default item delegate (item rendering style) here
	// but we customize the height depending on the view width during update cycle
	listView := list.New(entryListItems, models.NewItemDelegate(0), 0, 0)
	listView.SetShowHelp(false)
	listView.SetShowTitle(false)
	listView.SetFilteringEnabled(false)
	listView.SetStatusBarItemName("entry", "entries")
	listView.DisableQuitKeybindings()

	// Sidebar content
	randomEntryPtr := getNewRandomEntry(database, context)

	defaultBorderStyle := lipgloss.NewStyle().Border(lipgloss.RoundedBorder())
	return HomeView{
		width:        0,
		height:       0,
		sidebarWidth: 0,
		help:         help.New(),
		keymap: homeViewKeymap{
			quit: key.NewBinding(
				key.WithKeys("q"),
				key.WithHelp("q", "Quit"),
			),
			reroll: key.NewBinding(
				key.WithKeys("r"),
				key.WithHelp("r", "Reroll random entry"),
			),
		},
		database:  database,
		dbContext: context,

		borderStyle:   defaultBorderStyle,
		mainAreaStyle: defaultBorderStyle,
		sidebarStyle:  defaultBorderStyle,

		entriesList: listView,
		randomEntry: randomEntryPtr,
	}
}

func (v HomeView) Init() tea.Cmd {
	return nil
}

func (v HomeView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	// Determine the type of update message from Bubbletea using switch case
	switch msg := msg.(type) {

	// Note that in Go, cases are not indented within the switch block

	// On Window Size update
	case tea.WindowSizeMsg:
		// Expand the home view to fit the entire window
		v.width = msg.Width
		v.height = msg.Height - helpHeight
		// Recalculate sidebar width
		v.sidebarWidth = v.width * sidebarProportion / 100
		// Update styles
		borderSizeX := v.borderStyle.GetBorderLeftSize() + v.borderStyle.GetBorderRightSize()
		borderSizeY := v.borderStyle.GetBorderTopSize() + v.borderStyle.GetBorderBottomSize()
		mainAreaWidth := v.width - v.sidebarWidth - borderSizeX
		v.mainAreaStyle = v.borderStyle.
			Width(mainAreaWidth).
			Height(v.height - borderSizeY)
		v.sidebarStyle = v.borderStyle.
			Width(v.sidebarWidth - borderSizeX).
			Height(v.height - borderSizeY)
		// Update list size (fill main area)
		v.entriesList.SetSize(v.mainAreaStyle.GetWidth(), v.mainAreaStyle.GetHeight())
		// Update list item size
		v.entriesList.SetDelegate(models.NewItemDelegate(mainAreaWidth))

	// On keypress
	case tea.KeyMsg:
		switch {
		// q to quit
		case key.Matches(msg, v.keymap.quit):
			return v, tea.Quit
		// r to reroll random entry
		case key.Matches(msg, v.keymap.reroll):
			v.randomEntry = getNewRandomEntry(v.database, v.dbContext)
			return v, nil
		}
	}

	// Update child list view
	var listCmd tea.Cmd
	v.entriesList, listCmd = v.entriesList.Update(msg)

	// Return the updated model and no command
	return v, listCmd
}

func (v HomeView) View() string {

	helpView := v.help.ShortHelpView(append(v.entriesList.ShortHelp(), v.keymap.quit))

	mainAreaView := lipgloss.JoinHorizontal(lipgloss.Top,
		v.mainAreaStyle.Render(v.entriesList.View()),
		v.sidebarStyle.Render(v.getSidebarContent()))

	ui := lipgloss.JoinVertical(lipgloss.Left,
		mainAreaView,
		helpView,
	)

	// Return the UI for rendering
	return ui
}

func (v HomeView) getSidebarContent() string {

	sidebarContentWidth := v.sidebarWidth - v.sidebarStyle.GetBorderLeftSize() - v.sidebarStyle.GetBorderRightSize()
	headerStyle := lipgloss.NewStyle().Bold(true).Width(sidebarContentWidth)
	contentStyle := lipgloss.NewStyle().Italic(true).Width(sidebarContentWidth)

	if v.randomEntry != nil {
		title := headerStyle.Render(fmt.Sprintf("%s...", v.randomEntry.GetRelativeDateString()))
		content := contentStyle.Render(v.randomEntry.Content)
		return lipgloss.JoinVertical(lipgloss.Left, title, "", content)
	} else {
		return "Time waits for no one..."
	}

}

func getNewRandomEntry(database *db.Queries, context context.Context) *db.Entry {
	randomEntry, err := database.GetRandomEntry(context)
	var randomEntryPtr *db.Entry
	if err == nil {
		randomEntryPtr = &randomEntry
	}
	return randomEntryPtr
}
